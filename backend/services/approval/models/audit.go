package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AuditActionApproved = "approved"
	AuditActionDenied   = "denied"
)

type ApprovalAudit struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookingID uuid.UUID `gorm:"type:uuid;index"`
	StaffID   uuid.UUID `gorm:"type:uuid"`
	Action    string    `gorm:"type:varchar(20)"`
	Reason    string    `gorm:"type:text"`
	CreatedAt time.Time
}
