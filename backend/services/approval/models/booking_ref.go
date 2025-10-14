package models

import (
	"time"

	"github.com/google/uuid"
)

// Booking is a lightweight mirror of the bookings table from booking service.
// We only include the fields we need to read/update.
type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	RoomID    uuid.UUID `gorm:"type:uuid"`
	StartTime time.Time
	EndTime   time.Time
	Status    string `gorm:"type:varchar(20)"` // expected: pending, approved, denied, transferred, etc.
}
