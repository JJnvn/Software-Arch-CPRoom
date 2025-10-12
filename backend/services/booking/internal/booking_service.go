package internal

import (
	"context"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingService struct {
	repo *BookingRepository
	pb.UnimplementedBookingServiceServer
}

func NewBookingService(repo *BookingRepository) *BookingService {
	return &BookingService{repo: repo}
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

	if err := s.repo.CancelBooking(id); err != nil {
		return &pb.CancelBookingResponse{Success: false}, err
	}
	return &pb.CancelBookingResponse{Success: true}, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *pb.UpdateBookingRequest) (*pb.UpdateBookingResponse, error) {
	id, err := uuid.Parse(req.BookingId)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateBooking(id, req.NewStart.AsTime(), req.NewEnd.AsTime())
	if err != nil {
		return &pb.UpdateBookingResponse{Success: false}, err
	}
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

	err = s.repo.TransferBooking(id, newOwner)
	if err != nil {
		return &pb.TransferBookingResponse{Success: false}, err
	}
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
