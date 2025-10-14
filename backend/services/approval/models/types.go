package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusDenied   Status = "DENIED"
)

// GORM models

type Booking struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID            string    `gorm:"not null"`
	UserID            string    `gorm:"not null"`
	Start             time.Time `gorm:"not null"`
	End               time.Time `gorm:"not null"`
	Status            Status    `gorm:"type:text;not null;default:PENDING"`
	CurrentApproverID string    `gorm:"type:text"`
	RequestedAt       time.Time `gorm:"not null;default:now()"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Booking) TableName() string { return "bookings" }

type AuditEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	BookingID uuid.UUID `gorm:"type:uuid;index;not null"`
	StaffID   string    `gorm:"type:text;not null"` // actor ("system" allowed)
	Action    string    `gorm:"type:text;not null"` // approve|deny|reassign|assign
	Reason    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (AuditEvent) TableName() string { return "audit_events" }
