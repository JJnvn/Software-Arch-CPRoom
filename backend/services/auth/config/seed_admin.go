package config

import (
	"errors"
	"log"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminEmail    = "admin@admin.com"
	defaultAdminPassword = "Secured1"
	defaultAdminName     = "System Admin"
)

func SeedAdmin(db *gorm.DB) {
	var admin models.User
	err := db.Where("email = ?", defaultAdminEmail).First(&admin).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			log.Printf("Failed to hash admin password: %v", hashErr)
			return
		}

		admin = models.User{
			Name:     defaultAdminName,
			Email:    defaultAdminEmail,
			Password: string(hashed),
			Role:     models.ADMIN,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
			return
		}

		log.Println("Default admin created: admin@admin.com / Secured1")
		return

	case err != nil:
		log.Printf("Failed to load admin user: %v", err)
		return
	}

	updates := map[string]any{}

	if admin.Name == "" {
		updates["name"] = defaultAdminName
	}

	if admin.Role != models.ADMIN {
		updates["role"] = models.ADMIN
	}

	if bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(defaultAdminPassword)) != nil {
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			log.Printf("Failed to hash admin password: %v", hashErr)
			return
		}
		updates["password"] = string(hashed)
	}

	if len(updates) == 0 {
		log.Println("Admin already exists — skipping seeding")
		return
	}

	if err := db.Model(&admin).Updates(updates).Error; err != nil {
		log.Printf("Failed to update admin user: %v", err)
		return
	}

	log.Println("Default admin ensured: admin@admin.com / Secured1")
}
