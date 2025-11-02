package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusDenied    = "denied"
)

type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID    uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	StartTime time.Time
	EndTime   time.Time
	Status    string
}
