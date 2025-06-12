package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	"gorm.io/gorm"
)

func SeedProductForms(db *gorm.DB) {
	productForms := []models.ProductForm{
		{
			ID:           1,
			Form:         "Physical",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748554985/physical_ltedby.jpg",
			CloudImageID: 2,
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			Form:         "Digital",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748554986/digital_xmbpnw.jpg",
			CloudImageID: 2,
			CreatedAt:    time.Now(),
		},
		{
			ID:           3,
			Form:         "Service",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748555059/service_tx3tup.jpg",
			CloudImageID: 2,
			CreatedAt:    time.Now(),
		},
	}

	log.Printf("[DEBUG] Menemukan %d product forms untuk seeding", len(productForms))

	for _, pf := range productForms {
		pf.Slug = utils.SlugifyText(pf.Form)
		if err := db.FirstOrCreate(&pf, models.ProductForm{ID: pf.ID}).Error; err != nil {
			log.Printf("❌ Gagal seeding product form ID %d: %v\n", pf.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  product forms seeded")
}
