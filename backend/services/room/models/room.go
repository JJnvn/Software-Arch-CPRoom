package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// StringList persists a slice of strings as JSON in the database.
type StringList []string

// Value converts the slice to JSON bytes before storing.
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(s))
}

// Scan restores the slice from JSON stored in the database.
func (s *StringList) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("models.StringList: unsupported scan type %T", value)
	}

	if len(data) == 0 {
		*s = nil
		return nil
	}

	return json.Unmarshal(data, s)
}

type Room struct {
	ID       uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name     string     `json:"name"`
	Capacity int        `json:"capacity"`
	Features StringList `gorm:"type:jsonb" json:"features"`
}
