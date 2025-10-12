package models

type User struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name       string
	Email      string `gorm:"uniqueIndex"`
	Password   string
	Language   string
	Bookings   []Booking   `gorm:"foreignKey:UserID"`
	Preference Preferences `gorm:"foreignKey:UserID"`
}

type Booking struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string
	RoomName  string
	StartTime string
	EndTime   string
	Status    string
}

type Preferences struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	UserID           string `gorm:"uniqueIndex"`
	NotificationType string
	Language         string
}
