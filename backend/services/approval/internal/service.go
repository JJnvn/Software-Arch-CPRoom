package internal

import (
	"context"
	"errors"
	"log"

	events "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/events"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApprovalService struct {
    repo      *ApprovalRepository
    publisher events.Publisher
    pb.UnimplementedApprovalServiceServer
}

func NewApprovalService(repo *ApprovalRepository, publisher events.Publisher) *ApprovalService {
    return &ApprovalService{repo: repo, publisher: publisher}
}

func (s *ApprovalService) GetRoomNames(ids []string) (map[string]string, error) {
    return s.repo.GetRoomNames(ids)
}

func (s *ApprovalService) GetUserNames(ids []string) (map[string]string, error) {
    return s.repo.GetUserNames(ids)
}

func (s *ApprovalService) ListApproved() ([]PendingBooking, error) {
    return s.repo.ListApprovedBookings()
}

func (s *ApprovalService) publishBookingEvent(ctx context.Context, booking *models.Booking, event string, metadata map[string]any) {
	if s.publisher == nil || booking == nil {
		return
	}
	payload := events.BookingEvent{
		Event:     event,
		BookingID: booking.ID.String(),
		UserID:    booking.UserID.String(),
		RoomID:    booking.RoomID.String(),
		Status:    booking.Status,
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
		Metadata:  metadata,
	}
	if err := s.publisher.PublishBookingEvent(ctx, payload); err != nil {
		log.Printf("failed to publish approval event %s: %v", event, err)
	}
}

func (s *ApprovalService) ListPending(ctx context.Context, req *pb.ListPendingRequest) (*pb.ListPendingResponse, error) {
	bookings, err := s.repo.ListPendingBookings()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list pending bookings: %v", err)
	}

	resp := &pb.ListPendingResponse{}
	for _, b := range bookings {
		resp.Pending = append(resp.Pending, &pb.PendingBooking{
			BookingId: b.ID.String(),
			RoomId:    b.RoomID.String(),
			UserId:    b.UserID.String(),
			Start:     timestamppb.New(b.StartTime),
			End:       timestamppb.New(b.EndTime),
		})
	}
	return resp, nil
}

func (s *ApprovalService) ApproveBooking(ctx context.Context, req *pb.ApproveRequest) (*pb.ApproveResponse, error) {
	bookingID, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid booking_id")
	}

	staffID, err := uuid.Parse(req.GetStaffId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff_id")
	}

	booking, err := s.repo.ApproveBooking(bookingID, staffID)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingNotFound):
			return nil, status.Error(codes.NotFound, "booking not found")
		case errors.Is(err, ErrAlreadyProcessed):
			return nil, status.Error(codes.FailedPrecondition, "booking already processed")
		case errors.Is(err, ErrNoStatusChange):
			return &pb.ApproveResponse{Success: true}, nil
		default:
			return nil, status.Errorf(codes.Internal, "failed to approve booking: %v", err)
		}
	}

	s.publishBookingEvent(ctx, booking, events.BookingApprovedEvent, map[string]any{
		"staff_id": staffID.String(),
	})

	return &pb.ApproveResponse{Success: true}, nil
}

func (s *ApprovalService) DenyBooking(ctx context.Context, req *pb.DenyRequest) (*pb.DenyResponse, error) {
	bookingID, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid booking_id")
	}

	staffID, err := uuid.Parse(req.GetStaffId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid staff_id")
	}

	reason := req.GetReason()
	if reason == "" {
		return nil, status.Error(codes.InvalidArgument, "reason is required when denying a booking")
	}

	booking, err := s.repo.DenyBooking(bookingID, staffID, reason)
	if err != nil {
		switch {
		case errors.Is(err, ErrBookingNotFound):
			return nil, status.Error(codes.NotFound, "booking not found")
		case errors.Is(err, ErrAlreadyProcessed):
			return nil, status.Error(codes.FailedPrecondition, "booking already processed")
		case errors.Is(err, ErrNoStatusChange):
			return &pb.DenyResponse{Success: true}, nil
		default:
			return nil, status.Errorf(codes.Internal, "failed to deny booking: %v", err)
		}
	}

	s.publishBookingEvent(ctx, booking, events.BookingDeniedEvent, map[string]any{
		"staff_id": staffID.String(),
		"reason":   reason,
	})

	return &pb.DenyResponse{Success: true}, nil
}

func (s *ApprovalService) GetAuditTrail(ctx context.Context, req *pb.GetAuditTrailRequest) (*pb.AuditTrailResponse, error) {
	bookingID, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid booking_id")
	}

	events, err := s.repo.GetAuditTrail(bookingID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch audit trail: %v", err)
	}

	resp := &pb.AuditTrailResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, &pb.AuditEvent{
			EventId:   e.ID.String(),
			BookingId: e.BookingID.String(),
			StaffId:   e.StaffID.String(),
			Action:    e.Action,
			Reason:    e.Reason,
			CreatedAt: timestamppb.New(e.CreatedAt),
		})
	}

	return resp, nil
}
