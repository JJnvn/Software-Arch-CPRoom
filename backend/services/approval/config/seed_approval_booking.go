package config

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedApprovalBookings(db *gorm.DB) {
	if db == nil {
		log.Println("SeedApprovalBookings skipped: database handle is nil")
		return
	}

	// Fixed user ID (replace with your seeded admin/user if needed)
	userID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"

	// Fixed rooms (must match rooms seeded earlier)
	roomIDs := []string{
		"11111111-1111-1111-1111-111111111111",
		"22222222-2222-2222-2222-222222222222",
		"33333333-3333-3333-3333-333333333333",
	}

	// Booking times
	now := time.Now()
	start1 := now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	end1 := now.Add(26 * time.Hour).Format("2006-01-02 15:04:05")

	start2 := now.Add(48 * time.Hour).Format("2006-01-02 15:04:05")
	end2 := now.Add(50 * time.Hour).Format("2006-01-02 15:04:05")

	start3 := now.Add(72 * time.Hour).Format("2006-01-02 15:04:05")
	end3 := now.Add(74 * time.Hour).Format("2006-01-02 15:04:05")

	queries := []string{
		// Confirmed booking
		`INSERT INTO bookings (id, user_id, room_id, start_time, end_time, status, created_at, updated_at)
		 VALUES ('aaaa1111-1111-1111-1111-111111111111', '` + userID + `', '` + roomIDs[0] + `', '` + start1 + `', '` + end1 + `', 'confirmed', NOW(), NOW())
		 ON CONFLICT (id) DO NOTHING;`,

		// Pending bookings
		`INSERT INTO bookings (id, user_id, room_id, start_time, end_time, status, created_at, updated_at)
		 VALUES ('aaaa2222-2222-2222-2222-222222222222', '` + userID + `', '` + roomIDs[1] + `', '` + start2 + `', '` + end2 + `', 'pending', NOW(), NOW())
		 ON CONFLICT (id) DO NOTHING;`,

		`INSERT INTO bookings (id, user_id, room_id, start_time, end_time, status, created_at, updated_at)
		 VALUES ('aaaa3333-3333-3333-3333-333333333333', '` + userID + `', '` + roomIDs[2] + `', '` + start3 + `', '` + end3 + `', 'pending', NOW(), NOW())
		 ON CONFLICT (id) DO NOTHING;`,
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Printf("Failed to insert approval booking: %v", err)
		}
	}

	log.Println("Seeded 3 approval bookings successfully (or already existed).")
}
