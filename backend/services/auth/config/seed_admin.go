package config

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminID       = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	defaultAdminEmail    = "admin@admin.com"
	defaultAdminPassword = "Secured1"
	defaultAdminName     = "System Admin"
	defaultAdminRole     = "admin"
)

func SeedAdmin(db *gorm.DB) {
	if db == nil {
		log.Println("SeedAdmin skipped: database handle is nil")
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	// Raw SQL insert with fixed ID
	sql := `
		INSERT INTO users (id, name, email, password, role)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    email = EXCLUDED.email,
		    password = EXCLUDED.password,
		    role = EXCLUDED.role;
	`

	if err := db.Exec(sql, defaultAdminID, defaultAdminName, defaultAdminEmail, string(hashed), defaultAdminRole).Error; err != nil {
		log.Printf("Failed to seed default admin: %v", err)
		return
	}

	log.Println("Default admin seeded or updated with fixed ID: admin@admin.com / Secured1")
}
