package internal

import (
	"context"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	pb "github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApprovalService struct {
	repo *ApprovalRepository
	pb.UnimplementedApprovalServiceServer
}

func NewApprovalService(repo *ApprovalRepository) *ApprovalService {
	return &ApprovalService{repo: repo}
}

func (s *ApprovalService) ListPending(ctx context.Context, req *pb.ListPendingRequest) (*pb.ListPendingResponse, error) {
	staffID, err := uuid.Parse(req.StaffId)
	if err != nil { return nil, err }
	rows, err := s.repo.ListPendingByStaff(staffID)
	if err != nil { return nil, err }
	resp := &pb.ListPendingResponse{}
	for _, r := range rows {
		resp.Pending = append(resp.Pending, &pb.PendingBooking{
			BookingId: r.BookingID.String(),
			RoomId:    r.RoomID.String(),
			UserId:    r.UserID.String(),
			Start:     timestamppb.New(r.Start),
			End:       timestamppb.New(r.End),
		})
	}
	return resp, nil
}

func (s *ApprovalService) ApproveBooking(ctx context.Context, req *pb.ApproveRequest) (*pb.ApproveResponse, error) {
	bid, err := uuid.Parse(req.BookingId)
	if err != nil { return nil, err }
	sid, err := uuid.Parse(req.StaffId)
	if err != nil { return nil, err }
	if err := s.repo.Approve(bid, sid); err != nil { return &pb.ApproveResponse{Success: false}, err }
	_ = s.repo.AddAudit(bid, sid, "approve", "")
	_ = s.repo.setBookingStatus(bid, "approved")
	return &pb.ApproveResponse{Success: true}, nil
}

func (s *ApprovalService) DenyBooking(ctx context.Context, req *pb.DenyRequest) (*pb.DenyResponse, error) {
	bid, err := uuid.Parse(req.BookingId)
	if err != nil { return nil, err }
	sid, err := uuid.Parse(req.StaffId)
	if err != nil { return nil, err }
	if err := s.repo.Deny(bid, sid, req.Reason); err != nil { return &pb.DenyResponse{Success: false}, err }
	_ = s.repo.AddAudit(bid, sid, "deny", req.Reason)
	_ = s.repo.setBookingStatus(bid, "denied")
	return &pb.DenyResponse{Success: true}, nil
}

func (s *ApprovalService) ReassignApprover(ctx context.Context, req *pb.ReassignRequest) (*pb.ReassignResponse, error) {
	bid, err := uuid.Parse(req.BookingId)
	if err != nil { return nil, err }
	nsid, err := uuid.Parse(req.NewStaffId)
	if err != nil { return nil, err }
	if err := s.repo.Reassign(bid, nsid); err != nil { return &pb.ReassignResponse{Success: false}, err }
	_ = s.repo.AddAudit(bid, nsid, "reassign", "")
	return &pb.ReassignResponse{Success: true}, nil
}

func (s *ApprovalService) GetAuditTrail(ctx context.Context, req *pb.GetAuditTrailRequest) (*pb.AuditTrailResponse, error) {
	bid, err := uuid.Parse(req.BookingId)
	if err != nil { return nil, err }
	evts, err := s.repo.GetAuditTrail(bid)
	if err != nil { return nil, err }
	resp := &pb.AuditTrailResponse{}
	for _, e := range evts {
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
