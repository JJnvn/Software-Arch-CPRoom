package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	notifier "github.com/JJnvn/Software-Arch-CPRoom/backend/libs/notifier"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	errNotPending = errors.New("booking not pending")
)

type ApprovalService struct {
	db       *gorm.DB
	repo     *Repository
	notifier *notifier.Client
	resolver *UserResolver
}

func NewApprovalService(db *gorm.DB, notifierClient *notifier.Client, resolver *UserResolver) *ApprovalService {
	return &ApprovalService{
		db:       db,
		repo:     NewRepository(),
		notifier: notifierClient,
		resolver: resolver,
	}
}

func (s *ApprovalService) ListPending() ([]models.Booking, error) {
	return s.repo.ListPending(s.db)
}

func (s *ApprovalService) Approve(ctx context.Context, id uuid.UUID, staffID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		booking, err := s.repo.GetForUpdate(tx, id)
		if err != nil {
			return err
		}
		if booking.Status != models.StatusPending {
			return errNotPending
		}

		now := time.Now()
		fields := map[string]any{
			"status":        models.StatusApproved,
			"approved_by":   staffID,
			"approved_at":   now,
			"denied_by":     nil,
			"denied_at":     nil,
			"denied_reason": nil,
		}
		if err := s.repo.Update(tx, id, fields); err != nil {
			return err
		}
		booking.Status = models.StatusApproved
		booking.ApprovedBy = &staffID
		booking.ApprovedAt = &now
		booking.DeniedBy = nil
		booking.DeniedAt = nil
		booking.DeniedReason = nil

		if err := s.sendApprovalNotification(ctx, booking); err != nil {
			return fmt.Errorf("notify: %w", err)
		}

		return nil
	})
}

func (s *ApprovalService) Deny(ctx context.Context, id uuid.UUID, staffID, reason string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		booking, err := s.repo.GetForUpdate(tx, id)
		if err != nil {
			return err
		}
		if booking.Status != models.StatusPending {
			return errNotPending
		}

		now := time.Now()
		fields := map[string]any{
			"status":        models.StatusDenied,
			"denied_by":     staffID,
			"denied_at":     now,
			"denied_reason": reason,
			"approved_by":   nil,
			"approved_at":   nil,
		}
		if err := s.repo.Update(tx, id, fields); err != nil {
			return err
		}
		booking.Status = models.StatusDenied
		booking.DeniedBy = &staffID
		booking.DeniedAt = &now
		booking.DeniedReason = &reason
		booking.ApprovedBy = nil
		booking.ApprovedAt = nil

		if err := s.sendDenialNotification(ctx, booking, reason); err != nil {
			return fmt.Errorf("notify: %w", err)
		}
		return nil
	})
}

func (s *ApprovalService) sendApprovalNotification(ctx context.Context, booking models.Booking) error {
	if s.notifier == nil || !s.notifier.Enabled() {
		return nil
	}

	email, err := s.resolveUserEmail(ctx, booking.UserID)
	if err != nil {
		return err
	}

	metadata := map[string]any{
		"booking_id": booking.ID,
		"room_id":    booking.RoomID,
		"start_time": booking.StartTime.Format(time.RFC3339),
		"end_time":   booking.EndTime.Format(time.RFC3339),
		"status":     models.StatusApproved,
		"email":      email,
	}
	message := fmt.Sprintf("Your booking for room %s has been approved.", booking.RoomID)

	if err := s.notifier.Send(ctx, booking.UserID, "booking_approved", message, metadata); err != nil {
		return err
	}

	reminderTime := booking.StartTime.Add(-30 * time.Minute)
	if reminderTime.After(time.Now()) {
		reminderMsg := fmt.Sprintf("Reminder: your room %s booking starts at %s", booking.RoomID, booking.StartTime.Format(time.RFC3339))
		if err := s.notifier.Schedule(ctx, booking.UserID, "booking_reminder", reminderMsg, reminderTime, metadata); err != nil {
			return err
		}
	}

	return nil
}

func (s *ApprovalService) sendDenialNotification(ctx context.Context, booking models.Booking, reason string) error {
	if s.notifier == nil || !s.notifier.Enabled() {
		return nil
	}
	email, err := s.resolveUserEmail(ctx, booking.UserID)
	if err != nil {
		return err
	}
	metadata := map[string]any{
		"booking_id": booking.ID,
		"room_id":    booking.RoomID,
		"start_time": booking.StartTime.Format(time.RFC3339),
		"end_time":   booking.EndTime.Format(time.RFC3339),
		"status":     models.StatusDenied,
		"reason":     reason,
		"email":      email,
	}
	message := fmt.Sprintf("Your booking for room %s has been denied.", booking.RoomID)
	return s.notifier.Send(ctx, booking.UserID, "booking_denied", message, metadata)
}

func (s *ApprovalService) resolveUserEmail(ctx context.Context, userID string) (string, error) {
	if s.resolver == nil {
		return "", errors.New("user resolver not configured")
	}
	return s.resolver.ResolveEmail(ctx, userID)
}
