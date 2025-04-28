package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedDiscountSponsors(db *gorm.DB) {

	discountSponsors := []models.DiscountSponsor{
		{ID: 1, RoleID: 2, CreatedAt: time.Now()},
		{ID: 2, RoleID: 4, CreatedAt: time.Now()},
	}

	for _, discountSponsor := range discountSponsors {
		if err := db.FirstOrCreate(&discountSponsor, models.DiscountSponsor{ID: discountSponsor.ID}).Error; err != nil {
			log.Printf("Gagal seeding DiscountSponsor ID %d: %v\n", discountSponsor.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  discount sponsors seeded")
}
