package internal

import (
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/room/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*models.Room, error)
	List() ([]models.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *roomRepository {
	db.AutoMigrate(&models.Room{})
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *roomRepository) GetByID(id uuid.UUID) (*models.Room, error) {
	var room models.Room
	err := r.db.First(&room, id).Error
	return &room, err
}

func (r *roomRepository) List() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Find(&rooms).Error
	return rooms, err
}
