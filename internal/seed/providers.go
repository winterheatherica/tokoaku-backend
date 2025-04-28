package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedProviders(db *gorm.DB) {
	providers := []models.Provider{
		{ID: 1, Name: "Email & Password (Firebase)", CreatedAt: time.Now()},
		{ID: 2, Name: "Google (Firebase)", CreatedAt: time.Now()},
		{ID: 11, Name: "Cloudinary", CreatedAt: time.Now()},
		{ID: 12, Name: "Upstash (Redis)", CreatedAt: time.Now()},
	}

	for _, p := range providers {
		if err := db.FirstOrCreate(&p, models.Provider{ID: p.ID}).Error; err != nil {
			log.Printf("Gagal seeding provider ID %d: %v\n", p.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  providers seeded")
}
