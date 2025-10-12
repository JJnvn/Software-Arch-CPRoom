package internal

import (
	"time"

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

func (r *BookingRepository) CancelBooking(id uuid.UUID) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", "cancelled").Error
}

func (r *BookingRepository) UpdateBooking(id uuid.UUID, newStart, newEnd time.Time) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Updates(map[string]interface{}{
		"start_time": newStart,
		"end_time":   newEnd,
	}).Error
}

func (r *BookingRepository) TransferBooking(id, newOwner uuid.UUID) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("user_id", newOwner).Error
}

func (r *BookingRepository) GetRoomSchedule(roomID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("room_id = ? AND status = ?", roomID, "active").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) AdminListBookings(roomID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("room_id = ?", roomID).Find(&bookings).Error
	return bookings, err
}
