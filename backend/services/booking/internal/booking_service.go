package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingService struct {
	repo     *BookingRepository
	notifier *Notifier
	pb.UnimplementedBookingServiceServer
}

func NewBookingService(repo *BookingRepository, notifier *Notifier) *BookingService {
	return &BookingService{repo: repo, notifier: notifier}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	id := uuid.New()
	booking := &models.Booking{
		ID:        id,
		UserID:    uuid.MustParse(req.UserId),
		RoomID:    uuid.MustParse(req.RoomId),
		StartTime: req.Start.AsTime(),
		EndTime:   req.End.AsTime(),
		Status:    "active",
	}

	if err := s.repo.Create(booking); err != nil {
		return nil, err
	}

	s.notifyBookingCreated(ctx, booking)

	return &pb.CreateBookingResponse{
		BookingId: booking.ID.String(),
		UserId:    booking.UserID.String(),
		RoomId:    booking.RoomID.String(),
		Start:     req.Start,
		End:       req.End,
	}, nil
}

func (s *BookingService) GetRoomSchedule(ctx context.Context, req *pb.GetRoomScheduleRequest) (*pb.RoomScheduleResponse, error) {
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		return nil, err
	}

	bookings, err := s.repo.GetRoomSchedule(roomID)
	if err != nil {
		return nil, err
	}

	resp := &pb.RoomScheduleResponse{RoomId: req.RoomId}
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

func (s *BookingService) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	id, err := uuid.Parse(req.BookingId)
	if err != nil {
		return nil, err
	}

	booking, err := s.repo.FindByID(id)
	if err != nil {
		return &pb.CancelBookingResponse{Success: false}, err
	}

	if err := s.repo.CancelBooking(id); err != nil {
		return &pb.CancelBookingResponse{Success: false}, err
	}
	booking.Status = "cancelled"
	s.notifyBookingCancelled(ctx, booking)
	return &pb.CancelBookingResponse{Success: true}, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.UpdateBookingResponse, error) {
	id, err := uuid.Parse(req.BookingId)
	if err != nil {
		return nil, err
	}

	booking, err := s.repo.FindByID(id)
	if err != nil {
		return &pb.UpdateBookingResponse{Success: false}, err
	}

	newStart := req.NewStart.AsTime()
	newEnd := req.NewEnd.AsTime()

	err = s.repo.UpdateBooking(id, newStart, newEnd)
	if err != nil {
		return &pb.UpdateBookingResponse{Success: false}, err
	}
	booking.StartTime = newStart
	booking.EndTime = newEnd
	s.notifyBookingUpdated(ctx, booking)
	return &pb.UpdateBookingResponse{Success: true}, nil
}

func (s *BookingService) TransferBooking(ctx context.Context, req *pb.TransferBookingRequest) (*pb.TransferBookingResponse, error) {
	id, err := uuid.Parse(req.BookingId)
	if err != nil {
		return nil, err
	}

	newOwner, err := uuid.Parse(req.NewOwnerId)
	if err != nil {
		return nil, err
	}

	booking, err := s.repo.FindByID(id)
	if err != nil {
		return &pb.TransferBookingResponse{Success: false}, err
	}

	err = s.repo.TransferBooking(id, newOwner)
	if err != nil {
		return &pb.TransferBookingResponse{Success: false}, err
	}
	booking.UserID = newOwner
	s.notifyBookingTransferred(ctx, booking)
	return &pb.TransferBookingResponse{Success: true}, nil
}

func (s *BookingService) AdminListBookings(ctx context.Context, req *pb.AdminListBookingsRequest) (*pb.AdminListBookingsResponse, error) {
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		return nil, err
	}

	bookings, err := s.repo.AdminListBookings(roomID)
	if err != nil {
		return nil, err
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

func (s *BookingService) notifyBookingCreated(ctx context.Context, booking *models.Booking) {
	if s.notifier == nil {
		return
	}
	metadata := bookingMetadata(booking)
	message := fmt.Sprintf("Booking confirmed for room %s from %s to %s", booking.RoomID, formatTime(booking.StartTime), formatTime(booking.EndTime))
	if err := s.notifier.SendImmediate(ctx, booking.UserID.String(), "booking_created", message, metadata); err != nil {
		log.Printf("notification send failed: %v", err)
	}

	reminderTime := booking.StartTime.Add(-30 * time.Minute)
	if reminderTime.After(time.Now()) {
		reminderMsg := fmt.Sprintf("Reminder: room %s booking starts at %s", booking.RoomID, formatTime(booking.StartTime))
		if err := s.notifier.Schedule(ctx, booking.UserID.String(), "booking_reminder", reminderMsg, reminderTime, metadata); err != nil {
			log.Printf("schedule reminder failed: %v", err)
		}
	}
}

func (s *BookingService) notifyBookingUpdated(ctx context.Context, booking *models.Booking) {
	if s.notifier == nil {
		return
	}
	metadata := bookingMetadata(booking)
	message := fmt.Sprintf("Booking updated for room %s. New time: %s - %s", booking.RoomID, formatTime(booking.StartTime), formatTime(booking.EndTime))
	if err := s.notifier.SendImmediate(ctx, booking.UserID.String(), "booking_updated", message, metadata); err != nil {
		log.Printf("notification send failed: %v", err)
	}

	reminderTime := booking.StartTime.Add(-30 * time.Minute)
	if reminderTime.After(time.Now()) {
		reminderMsg := fmt.Sprintf("Reminder: room %s booking now starts at %s", booking.RoomID, formatTime(booking.StartTime))
		if err := s.notifier.Schedule(ctx, booking.UserID.String(), "booking_reminder", reminderMsg, reminderTime, metadata); err != nil {
			log.Printf("schedule reminder failed: %v", err)
		}
	}
}

func (s *BookingService) notifyBookingCancelled(ctx context.Context, booking *models.Booking) {
	if s.notifier == nil {
		return
	}
	metadata := bookingMetadata(booking)
	message := fmt.Sprintf("Booking for room %s on %s has been cancelled", booking.RoomID, formatTime(booking.StartTime))
	if err := s.notifier.SendImmediate(ctx, booking.UserID.String(), "booking_cancelled", message, metadata); err != nil {
		log.Printf("notification send failed: %v", err)
	}
}

func (s *BookingService) notifyBookingTransferred(ctx context.Context, booking *models.Booking) {
	if s.notifier == nil {
		return
	}
	metadata := bookingMetadata(booking)
	message := fmt.Sprintf("You have been assigned a booking for room %s on %s", booking.RoomID, formatTime(booking.StartTime))
	if err := s.notifier.SendImmediate(ctx, booking.UserID.String(), "booking_transferred", message, metadata); err != nil {
		log.Printf("notification send failed: %v", err)
	}
}

func bookingMetadata(booking *models.Booking) map[string]any {
	return map[string]any{
		"booking_id": booking.ID.String(),
		"user_id":    booking.UserID.String(),
		"room_id":    booking.RoomID.String(),
		"start_time": booking.StartTime.UTC().Format(time.RFC3339),
		"end_time":   booking.EndTime.UTC().Format(time.RFC3339),
		"status":     booking.Status,
	}
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
