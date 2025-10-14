package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound   = errors.New("booking not found")
	ErrNotPending = errors.New("booking is not pending")
)

// Cursor helpers: base64url("unix_nano|uuid")
func encodeToken(ts time.Time, id uuid.UUID) string {
	raw := fmt.Sprintf("%d|%s", ts.UnixNano(), id.String())
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}
func decodeToken(tok string) (*time.Time, *uuid.UUID, error) {
	if tok == "" {
		return nil, nil, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(tok)
	if err != nil {
		return nil, nil, fmt.Errorf("bad page_token: %w", err)
	}
	parts := strings.SplitN(string(b), "|", 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("bad page_token format")
	}
	var ns int64
	if _, err := fmt.Sscanf(parts[0], "%d", &ns); err != nil {
		return nil, nil, fmt.Errorf("bad page_token ts: %w", err)
	}
	id, err := uuid.Parse(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("bad page_token id: %w", err)
	}
	t := time.Unix(0, ns)
	return &t, &id, nil
}

type ApprovalRepo interface {
	ListPendingPage(db *gorm.DB, limit int, pageToken string) (rows []models.Booking, nextToken string, err error)
	CreatePending(db *gorm.DB, b models.Booking, assignedBy, reason string) (models.Booking, error)
	Approve(db *gorm.DB, bookingID uuid.UUID, staffID string) error
	Deny(db *gorm.DB, bookingID uuid.UUID, staffID, reason string) error
	Audit(db *gorm.DB, bookingID uuid.UUID) ([]models.AuditEvent, error)
}

type GormApprovalRepo struct{}

func NewGormApprovalRepo() *GormApprovalRepo { return &GormApprovalRepo{} }

// ------- List with pagination -------

func (r *GormApprovalRepo) ListPendingPage(db *gorm.DB, limit int, pageToken string) ([]models.Booking, string, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var cursorTs *time.Time
	var cursorID *uuid.UUID
	if pageToken != "" {
		t, id, err := decodeToken(pageToken)
		if err != nil {
			return nil, "", err
		}
		cursorTs, cursorID = t, id
	}

	q := db.Where("status = ?", models.StatusPending).
		Order("requested_at ASC").
		Order("id ASC")

	if cursorTs != nil && cursorID != nil {
		q = q.Where("(requested_at > ?) OR (requested_at = ? AND id > ?)", *cursorTs, *cursorTs, *cursorID)
	}

	var rows []models.Booking
	if err := q.Limit(limit + 1).Find(&rows).Error; err != nil {
		return nil, "", err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
		last := rows[len(rows)-1]
		return rows, encodeToken(last.RequestedAt, last.ID), nil
	}
	return rows, "", nil
}

// ------- Create -------

func (r *GormApprovalRepo) CreatePending(db *gorm.DB, b models.Booking, assignedBy, reason string) (models.Booking, error) {
	return b, db.Transaction(func(tx *gorm.DB) error {
		if b.ID == uuid.Nil {
			b.ID = uuid.New()
		}
		if b.Status == "" {
			b.Status = models.StatusPending
		}
		if b.RequestedAt.IsZero() {
			b.RequestedAt = time.Now()
		}

		if err := tx.Create(&b).Error; err != nil {
			return err
		}
		evt := models.AuditEvent{
			ID:        uuid.New(),
			BookingID: b.ID,
			StaffID:   assignedBy,
			Action:    "assign",
			Reason:    reason,
			CreatedAt: time.Now(),
		}
		return tx.Create(&evt).Error
	})
}

// ------- Decisions (any staff may approve/deny) -------

func (r *GormApprovalRepo) Approve(db *gorm.DB, bookingID uuid.UUID, staffID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var b models.Booking
		if err := tx.First(&b, "id = ?", bookingID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if b.Status != models.StatusPending {
			return ErrNotPending
		}

		if err := tx.Model(&models.Booking{}).
			Where("id = ? AND status = ?", bookingID, models.StatusPending).
			Update("status", models.StatusApproved).Error; err != nil {
			return err
		}
		evt := models.AuditEvent{
			ID:        uuid.New(),
			BookingID: bookingID,
			StaffID:   staffID,
			Action:    "approve",
			Reason:    "",
			CreatedAt: time.Now(),
		}
		return tx.Create(&evt).Error
	})
}

func (r *GormApprovalRepo) Deny(db *gorm.DB, bookingID uuid.UUID, staffID, reason string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var b models.Booking
		if err := tx.First(&b, "id = ?", bookingID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if b.Status != models.StatusPending {
			return ErrNotPending
		}

		if err := tx.Model(&models.Booking{}).
			Where("id = ? AND status = ?", bookingID, models.StatusPending).
			Update("status", models.StatusDenied).Error; err != nil {
			return err
		}
		evt := models.AuditEvent{
			ID:        uuid.New(),
			BookingID: bookingID,
			StaffID:   staffID,
			Action:    "deny",
			Reason:    reason,
			CreatedAt: time.Now(),
		}
		return tx.Create(&evt).Error
	})
}

// ------- Audit -------

func (r *GormApprovalRepo) Audit(db *gorm.DB, bookingID uuid.UUID) ([]models.AuditEvent, error) {
	var out []models.AuditEvent
	if err := db.Where("booking_id = ?", bookingID).
		Order("created_at ASC").Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
