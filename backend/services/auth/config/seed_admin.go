package config

import (
	"log"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultAdminEmail    = "admin@admin.com"
	defaultAdminPassword = "Secured1"
	defaultAdminName     = "System Admin"
)

func SeedAdmin(db *gorm.DB) {
	if db == nil {
		log.Println("SeedAdmin skipped: database handle is nil")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	admin := models.User{
		Name:     defaultAdminName,
		Email:    defaultAdminEmail,
		Password: string(hashed),
		Role:     models.ADMIN,
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&admin)

	if result.Error != nil {
		log.Printf("Failed to ensure default admin: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Println("Default admin created: admin@admin.com / Secured1")
		return
	}

	var existing models.User
	if err := db.Where("email = ?", defaultAdminEmail).First(&existing).Error; err != nil {
		log.Printf("Failed to load existing admin for update: %v", err)
		return
	}

	updates := map[string]any{}

	if existing.Name == "" {
		updates["name"] = defaultAdminName
	}
	if existing.Role != models.ADMIN {
		updates["role"] = models.ADMIN
	}
	if bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(defaultAdminPassword)) != nil {
		updates["password"] = string(hashed)
	}

	if len(updates) == 0 {
		log.Println("Default admin already present")
		return
	}

	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		log.Printf("Failed to refresh default admin: %v", err)
		return
	}

	log.Println("Default admin refreshed to expected defaults")
}
