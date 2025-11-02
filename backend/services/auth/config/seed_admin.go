package config

import (
	"log"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		log.Printf("Failed to check for existing admin: %v", err)
		return
	}

	if count > 0 {
		log.Println("Admin already exists — skipping seeding")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("Secured1"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	admin := &models.User{
		Name:     "System Admin",
		Email:    "admin@admin.com",
		Password: string(hashed),
		Role:     models.ADMIN,
	}

	if err := db.Create(admin).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return
	}

	log.Println("Default admin created: admin@admin.com / Secured1")
}
