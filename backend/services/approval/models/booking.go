package models

import "time"

const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusDenied    = "denied"
	StatusCancelled = "cancelled"
)

type Booking struct {
	ID           string     `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID       string     `gorm:"column:room_id" json:"room_id"`
	UserID       string     `gorm:"column:user_id" json:"user_id"`
	StartTime    time.Time  `gorm:"column:start_time" json:"start_time"`
	EndTime      time.Time  `gorm:"column:end_time" json:"end_time"`
	Status       string     `gorm:"column:status" json:"status"`
	ApprovedBy   *string    `gorm:"column:approved_by" json:"approved_by"`
	ApprovedAt   *time.Time `gorm:"column:approved_at" json:"approved_at"`
	DeniedBy     *string    `gorm:"column:denied_by" json:"denied_by"`
	DeniedAt     *time.Time `gorm:"column:denied_at" json:"denied_at"`
	DeniedReason *string    `gorm:"column:denied_reason" json:"denied_reason"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Booking) TableName() string {
	return "bookings"
}
