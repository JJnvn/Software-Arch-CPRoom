package internal

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/booking/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// ErrTimeSlotUnavailable is returned when a room is already booked (confirmed) for the requested window.
	ErrTimeSlotUnavailable = errors.New("room is not available for the requested time window")
	ErrRoomNotFound        = errors.New("room not found")
	ErrUserNotFound        = errors.New("user not found")
)

type BookingRepository struct {
	db *gorm.DB
}

type RoomSearchResult struct {
	ID       uuid.UUID
	Name     string
	Capacity int
	Features []string
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(b *models.Booking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var roomCount int64
		if err := tx.Table("rooms").Where("id = ?", b.RoomID).Count(&roomCount).Error; err != nil {
			return err
		}
		if roomCount == 0 {
			return ErrRoomNotFound
		}

		var userCount int64
		if err := tx.Table("users").Where("id = ?", b.UserID).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount == 0 {
			return ErrUserNotFound
		}

		var count int64
		if err := tx.Model(&models.Booking{}).
			Where("room_id = ?", b.RoomID).
			Where("status = ?", models.StatusConfirmed).
			Where("start_time < ? AND end_time > ?", b.EndTime, b.StartTime).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return ErrTimeSlotUnavailable
		}

		return tx.Create(b).Error
	})
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.First(&booking, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) UpdateStatus(id uuid.UUID, status string) error {
	result := r.db.Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BookingRepository) UpdateBookingTimes(id uuid.UUID, newStart, newEnd time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&booking, "id = ?", id).Error; err != nil {
			return err
		}

		var count int64
		if err := tx.Model(&models.Booking{}).
			Where("room_id = ?", booking.RoomID).
			Where("status = ?", models.StatusConfirmed).
			Where("id <> ?", booking.ID).
			Where("start_time < ? AND end_time > ?", newEnd, newStart).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return ErrTimeSlotUnavailable
		}

		booking.StartTime = newStart
		booking.EndTime = newEnd
		return tx.Save(&booking).Error
	})
}

func (r *BookingRepository) TransferBooking(id, newOwner uuid.UUID) error {
	result := r.db.Model(&models.Booking{}).
		Where("id = ?", id).
		Update("user_id", newOwner)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *BookingRepository) GetRoomSchedule(roomID uuid.UUID, date time.Time) ([]models.Booking, error) {
	var bookings []models.Booking

	// Set the start and end of the day for the given date
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("room_id = ?", roomID).
		Where("status IN ?", []string{
			models.StatusPending,
			models.StatusConfirmed,
			models.StatusExpired,
		}).
		Where("start_time < ? AND end_time > ?", endOfDay, startOfDay).
		Order("start_time ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) ListByUser(userID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("user_id = ?", userID).
		Order("start_time DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) AdminListBookings(roomID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("room_id = ?", roomID).
		Order("start_time ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) GetBookingsByRoom(roomID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("room_id = ?", roomID).
		Order("start_time DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) GetRoomBookingsByDate(roomID uuid.UUID, startOfDay, endOfDay time.Time) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("room_id = ?", roomID).
		Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay).
		Where("status IN ?", []string{models.StatusPending, models.StatusConfirmed}).
		Order("start_time ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) GetRoomName(roomID uuid.UUID) (string, error) {
	var roomName string
	err := r.db.Table("rooms").
		Select("name").
		Where("id = ?", roomID).
		Scan(&roomName).Error
	if err != nil {
		return "", err
	}
	return roomName, nil
}

func (r *BookingRepository) GetUserIDByEmail(email string) (uuid.UUID, error) {
	var userID string

	err := r.db.Table("users").
		Select("id").
		Where("email = ?", email).
		Scan(&userID).Error
	if err != nil {
		return uuid.Nil, err
	}
	if userID == "" {
		return uuid.Nil, ErrUserNotFound
	}
	return uuid.Parse(userID)
}

func (r *BookingRepository) SearchAvailableRooms(start, end time.Time, capacity, page, pageSize int) ([]RoomSearchResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	query := r.db.Table("rooms").
		Select("rooms.id, rooms.name, rooms.capacity, rooms.features")

	if capacity > 0 {
		query = query.Where("rooms.capacity >= ?", capacity)
	}

	if !start.IsZero() && !end.IsZero() {
		subQuery := r.db.Table("bookings").
			Select("1").
			Where("bookings.room_id = rooms.id").
			Where("bookings.status = ?", models.StatusConfirmed).
			Where("bookings.start_time < ? AND bookings.end_time > ?", end, start)
		query = query.Where("NOT EXISTS (?)", subQuery)
	}

	query = query.Order("rooms.capacity ASC, rooms.name ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	type roomSearchRow struct {
		ID       uuid.UUID
		Name     string
		Capacity int
		Features []byte
	}

	var rows []roomSearchRow
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	results := make([]RoomSearchResult, len(rows))
	for i, row := range rows {
		var featureSlice []string
		if len(row.Features) > 0 {
			if err := json.Unmarshal(row.Features, &featureSlice); err != nil {
				return nil, err
			}
		}

		results[i] = RoomSearchResult{
			ID:       row.ID,
			Name:     row.Name,
			Capacity: row.Capacity,
			Features: featureSlice,
		}
	}

	return results, nil
}
