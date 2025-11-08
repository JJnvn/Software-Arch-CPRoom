package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
	StatusExpired   = "expired"
	StatusCompleted = "completed"
	StatusDenied    = "denied"
)

type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	RoomID    uuid.UUID `gorm:"type:uuid;index:idx_room_time,priority:1" json:"room_id"`
	StartTime time.Time `gorm:"index:idx_room_time,priority:2" json:"start_time"`
	EndTime   time.Time `gorm:"index:idx_room_time,priority:3" json:"end_time"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
