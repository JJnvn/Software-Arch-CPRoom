package internal

import (
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalService struct {
	db   *gorm.DB
	repo *GormApprovalRepo
}

func NewApprovalService(db *gorm.DB, repo *GormApprovalRepo) *ApprovalService {
	return &ApprovalService{db: db, repo: repo}
}

// List
func (s *ApprovalService) ListPendingPage(limit int, token string) ([]models.Booking, string, error) {
	return s.repo.ListPendingPage(s.db, limit, token)
}

// Create
func (s *ApprovalService) CreatePending(roomID, userID string, start, end time.Time, assignedBy string) (models.Booking, error) {
	b := models.Booking{
		RoomID:      roomID,
		UserID:      userID,
		Start:       start,
		End:         end,
		Status:      models.StatusPending,
		RequestedAt: time.Now(),
	}
	return s.repo.CreatePending(s.db, b, assignedBy, "initial assignment")
}

// Decide
func (s *ApprovalService) Approve(id uuid.UUID, staffID string) error {
	return s.repo.Approve(s.db, id, staffID)
}
func (s *ApprovalService) Deny(id uuid.UUID, staffID, reason string) error {
	return s.repo.Deny(s.db, id, staffID, reason)
}

// Audit
func (s *ApprovalService) Audit(id uuid.UUID) ([]models.AuditEvent, error) {
	return s.repo.Audit(s.db, id)
}
