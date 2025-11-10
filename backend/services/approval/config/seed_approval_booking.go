package config

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedApprovalAudits(db *gorm.DB) {
	if db == nil {
		log.Println("SeedApprovalAudits skipped: database handle is nil")
		return
	}

	// Fixed staff ID (replace with your admin/staff UUID)
	staffID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"

	// Fixed booking IDs (from your seeded bookings)
	bookingIDs := []string{
		"aaaa1111-1111-1111-1111-111111111111", // confirmed booking
		"aaaa2222-2222-2222-2222-222222222222", // pending booking
		"aaaa3333-3333-3333-3333-333333333333", // pending booking
	}

	// Timestamps
	now := time.Now().Format("2006-01-02 15:04:05")

	queries := []string{
		// Approved audit
		`INSERT INTO approval_audits (id, booking_id, staff_id, action, reason, created_at)
	 VALUES ('11111111-1111-4111-8111-111111111111', '` + bookingIDs[0] + `', '` + staffID + `', 'approved', 'Automatically approved', '` + now + `')
	 ON CONFLICT (id) DO NOTHING;`,

		// Denied audits
		`INSERT INTO approval_audits (id, booking_id, staff_id, action, reason, created_at)
	 VALUES ('22222222-2222-4222-8222-222222222222', '` + bookingIDs[1] + `', '` + staffID + `', 'denied', 'Pending review', '` + now + `')
	 ON CONFLICT (id) DO NOTHING;`,

		`INSERT INTO approval_audits (id, booking_id, staff_id, action, reason, created_at)
	 VALUES ('33333333-3333-4333-8333-333333333333', '` + bookingIDs[2] + `', '` + staffID + `', 'denied', 'Pending review', '` + now + `')
	 ON CONFLICT (id) DO NOTHING;`,
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Printf("Failed to insert approval audit: %v", err)
		}
	}

	log.Println("Seeded 3 approval audits successfully (or already existed).")
}
