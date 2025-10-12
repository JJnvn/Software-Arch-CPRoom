package models

import "github.com/google/uuid"

type Room struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string    `json:"name"`
	Capacity int       `json:"capacity"`
}
