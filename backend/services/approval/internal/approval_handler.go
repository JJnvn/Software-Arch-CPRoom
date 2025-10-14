package internal

import (
	"context"
	"fmt"

	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApprovalHandler struct {
	pb.UnimplementedApprovalServiceServer
	svc *ApprovalService
}

func NewApprovalHandler(svc *ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{svc: svc}
}

func callerIDFromMD(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v := md.Get("x-user-id"); len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

// --- AddToList ---

func (h *ApprovalHandler) AddToList(ctx context.Context, req *pb.AddToListRequest) (*pb.AddToListResponse, error) {
	actor := callerIDFromMD(ctx)
	if actor == "" {
		actor = "system"
	}
	if req.GetRoomId() == "" || req.GetUserId() == "" {
		return nil, fmt.Errorf("room_id and user_id are required")
	}
	start := req.GetStart().AsTime()
	end := req.GetEnd().AsTime()
	if !end.After(start) {
		return nil, fmt.Errorf("end must be after start")
	}
	b, err := h.svc.CreatePending(req.GetRoomId(), req.GetUserId(), start, end, actor)
	if err != nil {
		return nil, err
	}
	return &pb.AddToListResponse{BookingId: b.ID.String()}, nil
}

// --- ListPending (paginated) ---

func (h *ApprovalHandler) ListPending(ctx context.Context, req *pb.ListPendingRequest) (*pb.ListPendingResponse, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	rows, next, err := h.svc.ListPendingPage(pageSize, req.GetPageToken())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListPendingResponse{Pending: make([]*pb.PendingBooking, 0, len(rows)), NextPageToken: next}
	for _, b := range rows {
		resp.Pending = append(resp.Pending, &pb.PendingBooking{
			BookingId: b.ID.String(),
			RoomId:    b.RoomID,
			UserId:    b.UserID,
			Start:     timestamppb.New(b.Start),
			End:       timestamppb.New(b.End),
			Status:    string(b.Status),
		})
	}
	return resp, nil
}

// --- Approve / Deny (any staff) ---

func (h *ApprovalHandler) ApproveBooking(ctx context.Context, req *pb.ApproveRequest) (*pb.ApproveResponse, error) {
	staffID := req.GetStaffId()
	if staffID == "" {
		staffID = callerIDFromMD(ctx)
	}
	bid, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking_id: %w", err)
	}
	if err := h.svc.Approve(bid, staffID); err != nil {
		switch err {
		case ErrNotPending:
			return &pb.ApproveResponse{Success: false}, nil
		default:
			return nil, err
		}
	}
	return &pb.ApproveResponse{Success: true}, nil
}

func (h *ApprovalHandler) DenyBooking(ctx context.Context, req *pb.DenyRequest) (*pb.DenyResponse, error) {
	staffID := req.GetStaffId()
	if staffID == "" {
		staffID = callerIDFromMD(ctx)
	}
	bid, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking_id: %w", err)
	}
	if err := h.svc.Deny(bid, staffID, req.GetReason()); err != nil {
		switch err {
		case ErrNotPending:
			return &pb.DenyResponse{Success: false}, nil
		default:
			return nil, err
		}
	}
	return &pb.DenyResponse{Success: true}, nil
}

// --- Audit ---

func (h *ApprovalHandler) GetAuditTrail(ctx context.Context, req *pb.GetAuditTrailRequest) (*pb.AuditTrailResponse, error) {
	bid, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, fmt.Errorf("invalid booking_id: %w", err)
	}
	evts, err := h.svc.Audit(bid)
	if err != nil {
		return nil, err
	}
	out := &pb.AuditTrailResponse{Events: make([]*pb.AuditEvent, 0, len(evts))}
	for _, e := range evts {
		out.Events = append(out.Events, &pb.AuditEvent{
			EventId:   e.ID.String(),
			BookingId: e.BookingID.String(),
			StaffId:   e.StaffID,
			Action:    e.Action,
			Reason:    e.Reason,
			CreatedAt: timestamppb.New(e.CreatedAt),
		})
	}
	return out, nil
}
