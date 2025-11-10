package internal

import (
	"context"
	"errors"
	"log"
	"time"

	events "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/events"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type BookingService struct {
	repo      *BookingRepository
	publisher events.Publisher
	pb.UnimplementedBookingServiceServer
}

func NewBookingService(repo *BookingRepository, publisher events.Publisher) *BookingService {
	return &BookingService{repo: repo, publisher: publisher}
}

func (s *BookingService) publishBookingEvent(ctx context.Context, booking *models.Booking, event string, metadata map[string]any) {
	if s.publisher == nil || booking == nil {
		return
	}

	// Fetch room name
	roomName, err := s.repo.GetRoomName(booking.RoomID)
	if err != nil {
		log.Printf("failed to fetch room name for room %s: %v", booking.RoomID, err)
		roomName = booking.RoomID.String() // Fallback to room ID
	}

	payload := events.BookingEvent{
		Event:     event,
		BookingID: booking.ID.String(),
		UserID:    booking.UserID.String(),
		RoomID:    booking.RoomID.String(),
		RoomName:  roomName,
		Status:    booking.Status,
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
		Metadata:  metadata,
	}
	if err := s.publisher.PublishBookingEvent(ctx, payload); err != nil {
		log.Printf("failed to publish booking event %s: %v", event, err)
	}
}

func (s *BookingService) SearchRooms(ctx context.Context, req *pb.SearchRoomsRequest) (*pb.SearchRoomsResponse, error) {
	start := req.GetStart().AsTime()
	end := req.GetEnd().AsTime()

	if start.IsZero() || end.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "start and end times are required")
	}

	if !start.Before(end) {
		return nil, status.Error(codes.InvalidArgument, "start time must be before end time")
	}

	rooms, err := s.repo.SearchAvailableRooms(
		start,
		end,
		int(req.Capacity),
		int(req.Page),
		int(req.PageSize),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to search rooms: %v", err)
	}

	resp := &pb.SearchRoomsResponse{}
	for _, room := range rooms {
		resp.Rooms = append(resp.Rooms, &pb.RoomInfo{
			RoomId:   room.ID.String(),
			Name:     room.Name,
			Capacity: int32(room.Capacity),
			Features: room.Features,
		})
	}

	return resp, nil
}

func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid room_id")
	}

	start := req.GetStart().AsTime()
	end := req.GetEnd().AsTime()

	if start.IsZero() || end.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "start and end times are required")
	}

	if !start.Before(end) {
		return nil, status.Error(codes.InvalidArgument, "start time must be before end time")
	}

	now := time.Now()
	if !start.After(now) {
		return nil, status.Error(codes.InvalidArgument, "start time must be in the future")
	}

	booking := &models.Booking{
		ID:        uuid.New(),
		UserID:    userID,
		RoomID:    roomID,
		StartTime: start,
		EndTime:   end,
		Status:    models.StatusPending,
	}

	if err := s.repo.Create(booking); err != nil {
		switch {
		case errors.Is(err, ErrTimeSlotUnavailable):
			return nil, status.Error(codes.FailedPrecondition, "room is not available for the requested time window")
		case errors.Is(err, ErrRoomNotFound):
			return nil, status.Error(codes.NotFound, "room not found")
		case errors.Is(err, ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to create booking: %v", err)
		}
	}

	s.publishBookingEvent(ctx, booking, events.BookingCreatedEvent, nil)

	return &pb.CreateBookingResponse{
		BookingId: booking.ID.String(),
		UserId:    booking.UserID.String(),
		RoomId:    booking.RoomID.String(),
		Start:     timestamppb.New(booking.StartTime),
		End:       timestamppb.New(booking.EndTime),
	}, nil
}

func (s *BookingService) GetRoomSchedule(ctx context.Context, req *pb.GetRoomScheduleRequest) (*pb.RoomScheduleResponse, error) {
	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid room_id")
	}

	// Parse the date string (format: YYYY-MM-DD)
	date := time.Now() // Default to today if not provided
	if req.GetDate() != "" {
		parsedDate, err := time.Parse("2006-01-02", req.GetDate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid date format, expected YYYY-MM-DD")
		}
		date = parsedDate
	}

	bookings, err := s.repo.GetRoomSchedule(roomID, date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load room schedule: %v", err)
	}

	resp := &pb.RoomScheduleResponse{RoomId: req.GetRoomId()}
	for _, b := range bookings {
		resp.Bookings = append(resp.Bookings, &pb.BookingSummary{
			BookingId: b.ID.String(),
			UserId:    b.UserID.String(),
			Start:     timestamppb.New(b.StartTime),
			End:       timestamppb.New(b.EndTime),
			Status:    b.Status,
		})
	}
	return resp, nil
}

func (s *BookingService) ListBookingsByUser(ctx context.Context, userID uuid.UUID) ([]models.Booking, error) {
	return s.repo.ListByUser(userID)
}

func (s *BookingService) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	id, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid booking_id")
	}

	booking, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to load booking: %v", err)
	}

	if booking.Status == models.StatusCancelled {
		return &pb.CancelBookingResponse{Success: true}, nil
	}

	now := time.Now()
	if !now.Before(booking.StartTime) {
		return nil, status.Error(codes.FailedPrecondition, "cannot cancel a booking that has already started")
	}

	if booking.Status == models.StatusExpired {
		return nil, status.Error(codes.FailedPrecondition, "booking already expired")
	}

	if err := s.repo.UpdateStatus(id, models.StatusCancelled); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel booking: %v", err)
	}

	booking.Status = models.StatusCancelled
	s.publishBookingEvent(ctx, booking, events.BookingCancelledEvent, nil)

	return &pb.CancelBookingResponse{Success: true}, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.UpdateBookingResponse, error) {
	id, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid booking_id")
	}

	newStart := req.GetNewStart().AsTime()
	newEnd := req.GetNewEnd().AsTime()

	if newStart.IsZero() || newEnd.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "new_start and new_end are required")
	}

	if !newStart.Before(newEnd) {
		return nil, status.Error(codes.InvalidArgument, "new_start must be before new_end")
	}

	now := time.Now()
	if !newStart.After(now) {
		return nil, status.Error(codes.InvalidArgument, "new_start must be in the future")
	}

	booking, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to load booking: %v", err)
	}

	switch booking.Status {
	case models.StatusCancelled, models.StatusExpired, models.StatusCompleted:
		return nil, status.Error(codes.FailedPrecondition, "booking cannot be modified in its current status")
	}

	if !now.Before(booking.StartTime) {
		return nil, status.Error(codes.FailedPrecondition, "cannot reschedule a booking that has already started")
	}

	if err := s.repo.UpdateBookingTimes(id, newStart, newEnd); err != nil {
		switch {
		case errors.Is(err, ErrTimeSlotUnavailable):
			return nil, status.Error(codes.FailedPrecondition, "room is not available for the requested time window")
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, "booking not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to update booking: %v", err)
		}
	}

	booking.StartTime = newStart
	booking.EndTime = newEnd
	s.publishBookingEvent(ctx, booking, events.BookingUpdatedEvent, nil)

	return &pb.UpdateBookingResponse{Success: true}, nil
}

func (s *BookingService) TransferBooking(ctx context.Context, req *pb.TransferBookingRequest) (*pb.TransferBookingResponse, error) {
	id, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid booking_id")
	}

	newOwner, err := uuid.Parse(req.GetNewOwnerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid new_owner_id")
	}

	booking, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to load booking: %v", err)
	}

	if booking.Status != models.StatusPending && booking.Status != models.StatusConfirmed {
		return nil, status.Error(codes.FailedPrecondition, "booking cannot be transferred in its current status")
	}

	now := time.Now()
	if !now.Before(booking.StartTime) {
		return nil, status.Error(codes.FailedPrecondition, "cannot transfer a booking that has already started")
	}

	if err := s.repo.TransferBooking(id, newOwner); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "booking not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to transfer booking: %v", err)
	}

	return &pb.TransferBookingResponse{Success: true}, nil
}

func (s *BookingService) AdminListBookings(ctx context.Context, req *pb.AdminListBookingsRequest) (*pb.AdminListBookingsResponse, error) {
	roomID, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid room_id")
	}

	bookings, err := s.repo.AdminListBookings(roomID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookings: %v", err)
	}

	resp := &pb.AdminListBookingsResponse{}
	for _, b := range bookings {
		resp.Bookings = append(resp.Bookings, &pb.BookingSummary{
			BookingId: b.ID.String(),
			UserId:    b.UserID.String(),
			Start:     timestamppb.New(b.StartTime),
			End:       timestamppb.New(b.EndTime),
			Status:    b.Status,
		})
	}
	return resp, nil
}
