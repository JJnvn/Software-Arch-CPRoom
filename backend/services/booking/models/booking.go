package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusDenied    = "denied"
	StatusCancelled = "cancelled"
)

type Booking struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	RoomID       uuid.UUID  `gorm:"type:uuid" json:"room_id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	Status       string     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ApprovedBy   *string    `gorm:"size:100" json:"approved_by"`
	ApprovedAt   *time.Time `json:"approved_at"`
	DeniedBy     *string    `gorm:"size:100" json:"denied_by"`
	DeniedAt     *time.Time `json:"denied_at"`
	DeniedReason *string    `gorm:"size:255" json:"denied_reason"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
