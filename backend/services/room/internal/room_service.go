package internal

import (
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/room/models"
	"github.com/google/uuid"
)

type RoomService interface {
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*models.Room, error)
	List() ([]models.Room, error)
}

type roomService struct {
	repo RoomRepository
}

func NewRoomService(repo *roomRepository) *roomService {
	return &roomService{repo: repo}
}

func (s *roomService) Create(room *models.Room) error {
	return s.repo.Create(room)
}

func (s *roomService) Update(room *models.Room) error {
	return s.repo.Update(room)
}

func (s *roomService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *roomService) GetByID(id uuid.UUID) (*models.Room, error) {
	return s.repo.GetByID(id)
}

func (s *roomService) List() ([]models.Room, error) {
	return s.repo.List()
}
