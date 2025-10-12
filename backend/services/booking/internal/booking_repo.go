package internal

import (
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(b *models.Booking) error {
	return r.db.Create(b).Error
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.First(&booking, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

// add more repo methods as needed (Update, Delete, ListByUser, etc.)
