package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedProductForms(db *gorm.DB) {
	productForms := []models.ProductForm{
		{ID: 1, Form: "Physical", CreatedAt: time.Now()},
		{ID: 2, Form: "Digital", CreatedAt: time.Now()},
		{ID: 3, Form: "Service", CreatedAt: time.Now()},
	}

	for _, pf := range productForms {
		if err := db.FirstOrCreate(&pf, models.ProductForm{ID: pf.ID}).Error; err != nil {
			log.Printf("Gagal seeding product_form ID %d: %v\n", pf.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  product forms seeded")
}
