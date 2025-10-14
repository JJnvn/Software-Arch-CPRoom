package internal

import (
	"errors"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/approval/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ListPending(db *gorm.DB) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := db.Where("status = ?", models.StatusPending).Order("created_at ASC").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *Repository) GetForUpdate(db *gorm.DB, id interface{}) (models.Booking, error) {
	var booking models.Booking
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, "id = ?", id).Error; err != nil {
		return models.Booking{}, err
	}
	return booking, nil
}

func (r *Repository) Update(db *gorm.DB, id interface{}, fields map[string]any) error {
	if len(fields) == 0 {
		return errors.New("no fields to update")
	}
	return db.Model(&models.Booking{}).Where("id = ?", id).Updates(fields).Error
}
