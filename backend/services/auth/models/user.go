package models

import "time"

type User struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string `gorm:"size:100;not null" json:"name"`
	Email     string `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

const TOKEN = "AUTH_TOKEN"
