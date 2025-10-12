package internal

import (
	"context"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/proto"
	"github.com/google/uuid"
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

// Add other gRPC methods (GetBooking, CancelBooking, UpdateBooking, etc.)
