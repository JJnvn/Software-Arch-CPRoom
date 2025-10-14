package internal

import (
	"errors"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

type PendingBookingRow struct {
	BookingID uuid.UUID
	RoomID    uuid.UUID
	UserID    uuid.UUID
	Start     time.Time
	End       time.Time
}

// ListPendingByStaff joins approvals with bookings for details.
func (r *ApprovalRepository) ListPendingByStaff(staffID uuid.UUID) ([]PendingBookingRow, error) {
	var rows []PendingBookingRow
	q := r.db.Table("approvals a").Select(
		"a.booking_id as booking_id, b.room_id, b.user_id, b.start_time as start, b.end_time as end").
		Joins("JOIN bookings b ON b.id = a.booking_id").
		Where("a.staff_id = ? AND a.status = ?", staffID, "pending")
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *ApprovalRepository) Approve(bookingID, staffID uuid.UUID) error {
	res := r.db.Model(&models.Approval{}).Where("booking_id = ? AND staff_id = ?", bookingID, staffID).Update("status", "approved")
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *ApprovalRepository) Deny(bookingID, staffID uuid.UUID, reason string) error {
	res := r.db.Model(&models.Approval{}).Where("booking_id = ? AND staff_id = ?", bookingID, staffID).Updates(map[string]any{"status": "denied", "reason": reason})
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *ApprovalRepository) Reassign(bookingID, newStaffID uuid.UUID) error {
	res := r.db.Model(&models.Approval{}).Where("booking_id = ?", bookingID).Update("staff_id", newStaffID)
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

func (r *ApprovalRepository) AddAudit(bookingID, staffID uuid.UUID, action, reason string) error {
	e := models.AuditEvent{ BookingID: bookingID, StaffID: staffID, Action: action, Reason: reason, CreatedAt: time.Now() }
	return r.db.Create(&e).Error
}

func (r *ApprovalRepository) GetAuditTrail(bookingID uuid.UUID) ([]models.AuditEvent, error) {
	var evts []models.AuditEvent
	if err := r.db.Where("booking_id = ?", bookingID).Order("created_at ASC").Find(&evts).Error; err != nil { return nil, err }
	return evts, nil
}

// Optional: keep booking table in sync for simple workflows. Non-fatal if not found.
func (r *ApprovalRepository) setBookingStatus(bookingID uuid.UUID, status string) error {
	res := r.db.Table("bookings").Where("id = ?", bookingID).Update("status", status)
	return res.Error
}

var ErrNotFound = errors.New("record not found")
