package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedProviders(db *gorm.DB) {
	providers := []models.Provider{
		{Name: "Email & Password (Firebase)", CreatedAt: time.Now()},
		{Name: "Google (Firebase)", CreatedAt: time.Now()},
		{Name: "Cloudinary", CreatedAt: time.Now()},
		{Name: "Upstash (Redis)", CreatedAt: time.Now()},
	}

	for _, p := range providers {
		if err := db.FirstOrCreate(&p, models.Provider{Name: p.Name}).Error; err != nil {
			log.Printf("Gagal seeding provider Name %s: %v\n", p.Name, err)
		}
	}

	log.Println("[SEEDER] ⚙️  Providers seeded")
}
