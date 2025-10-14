package models

import (
	"time"

	"github.com/google/uuid"
)

type Approval struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookingID uuid.UUID `gorm:"type:uuid;index"`
	StaffID   uuid.UUID `gorm:"type:uuid;index"`
	Status    string    `gorm:"type:varchar(20);default:'pending'"` // pending, approved, denied
	Reason    string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuditEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookingID uuid.UUID `gorm:"type:uuid;index"`
	StaffID   uuid.UUID `gorm:"type:uuid;index"`
	Action    string    `gorm:"type:varchar(20)"` // approve/deny/reassign
	Reason    string    `gorm:"type:text"`
	CreatedAt time.Time
}
