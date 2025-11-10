package config

import (
	"log"

	"gorm.io/gorm"
)

func SeedDefaultRooms(db *gorm.DB) {
	if db == nil {
		log.Println("SeedDefaultRooms skipped: database handle is nil")
		return
	}

	queries := []string{
		`INSERT INTO rooms (id, name, capacity, features)
		 VALUES ('11111111-1111-1111-1111-111111111111', 'Floor 1 Room 101', 10, '["Projector","Whiteboard","Air Conditioning","HDMI"]')
		 ON CONFLICT (id) DO NOTHING;`,

		`INSERT INTO rooms (id, name, capacity, features)
		 VALUES ('22222222-2222-2222-2222-222222222222', 'Floor 2 Room 203', 15, '["Projector","Whiteboard","HDMI"]')
		 ON CONFLICT (id) DO NOTHING;`,

		`INSERT INTO rooms (id, name, capacity, features)
		 VALUES ('33333333-3333-3333-3333-333333333333', 'Floor 3 Room 305', 8, '["Whiteboard","Air Conditioning"]')
		 ON CONFLICT (id) DO NOTHING;`,

		`INSERT INTO rooms (id, name, capacity, features)
		 VALUES ('44444444-4444-4444-4444-444444444444', 'Floor 5 Room 507', 20, '["Projector","Air Conditioning","HDMI"]')
		 ON CONFLICT (id) DO NOTHING;`,

		`INSERT INTO rooms (id, name, capacity, features)
		 VALUES ('55555555-5555-5555-5555-555555555555', 'Floor 9 Room 903', 12, '["Whiteboard","Projector","HDMI"]')
		 ON CONFLICT (id) DO NOTHING;`,
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Printf("Failed to insert room: %v", err)
		}
	}

	log.Println("Seeded 5 default rooms successfully (or already existed).")
}
