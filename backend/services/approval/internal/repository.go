package internal

import (
	"errors"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBookingNotFound  = errors.New("booking not found")
	ErrNoStatusChange   = errors.New("booking status unchanged")
	ErrAlreadyProcessed = errors.New("booking already processed")
)

type PendingBooking struct {
	ID        uuid.UUID
	RoomID    uuid.UUID
	UserID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

func (r *ApprovalRepository) ListPendingBookings() ([]PendingBooking, error) {
	var rows []PendingBooking
	err := r.db.Model(&models.Booking{}).
		Select("id, room_id, user_id, start_time, end_time").
		Where("status = ?", models.StatusPending).
		Order("start_time ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *ApprovalRepository) setBookingStatus(bookingID uuid.UUID, status string, staffID uuid.UUID, reason string, action string) (*models.Booking, error) {
	var updated *models.Booking
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, "id = ?", bookingID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookingNotFound
			}
			return err
		}

		if booking.Status != models.StatusPending {
			if booking.Status == status {
				return ErrNoStatusChange
			}
			return ErrAlreadyProcessed
		}

		now := time.Now().UTC()
		if err := tx.Model(&booking).Updates(map[string]any{
			"status":     status,
			"updated_at": now,
		}).Error; err != nil {
			return err
		}
		booking.Status = status
		booking.UpdatedAt = now

		event := &models.ApprovalAudit{
			ID:        uuid.New(),
			BookingID: bookingID,
			StaffID:   staffID,
			Action:    action,
			Reason:    reason,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		updated = &booking
		return nil
	})
	return updated, err
}

func (r *ApprovalRepository) ApproveBooking(bookingID, staffID uuid.UUID) (*models.Booking, error) {
	return r.setBookingStatus(bookingID, models.StatusConfirmed, staffID, "", models.AuditActionApproved)
}

func (r *ApprovalRepository) DenyBooking(bookingID, staffID uuid.UUID, reason string) (*models.Booking, error) {
	return r.setBookingStatus(bookingID, models.StatusDenied, staffID, reason, models.AuditActionDenied)
}

func (r *ApprovalRepository) GetAuditTrail(bookingID uuid.UUID) ([]models.ApprovalAudit, error) {
	var events []models.ApprovalAudit
	err := r.db.Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&events).Error
	return events, err
}
