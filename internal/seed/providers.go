package seed

import (
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedProviders(db *gorm.DB) {

	providers := []models.Provider{
		{ID: 1, Name: "password"},
		{ID: 2, Name: "google"},
		{ID: 3, Name: "github"},
	}

	for _, p := range providers {
		if err := db.FirstOrCreate(&p, models.Provider{ID: p.ID}).Error; err != nil {
			log.Printf("Gagal seeding provider ID %d: %v\n", p.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  Provider seeded")
}
