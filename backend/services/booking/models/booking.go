package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	RoomID    uuid.UUID `gorm:"type:uuid" json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
