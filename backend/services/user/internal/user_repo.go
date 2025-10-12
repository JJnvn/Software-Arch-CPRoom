package internal

import (
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetProfile(userID string) (*models.User, error)
	UpdateProfile(user *models.User) error
	GetBookingHistory(userID string) ([]models.Booking, error)
	GetPreferences(userID string) (*models.Preferences, error)
	UpdatePreferences(pref *models.Preferences) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetProfile(userID string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateProfile(user *models.User) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":     user.Name,
			"email":    user.Email,
			"language": user.Language,
		}).Error
}

func (r *userRepository) GetBookingHistory(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.db.Where("user_id = ?", userID).
		Order("start_time DESC").
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *userRepository) GetPreferences(userID string) (*models.Preferences, error) {
	var pref models.Preferences
	if err := r.db.First(&pref, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *userRepository) UpdatePreferences(pref *models.Preferences) error {
	return r.db.Model(&models.Preferences{}).
		Where("user_id = ?", pref.UserID).
		Updates(map[string]interface{}{
			"notification_type": pref.NotificationType,
			"language":          pref.Language,
		}).Error
}
