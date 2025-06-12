package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	"gorm.io/gorm"
)

func SeedCategoryLabels(db *gorm.DB) {
	categoryLabels := []models.CategoryLabel{
		{Name: "For", CreatedAt: time.Now()},
		{Name: "Opinion", CreatedAt: time.Now()},
		{Name: "Size", CreatedAt: time.Now()},
		{Name: "Product Age", CreatedAt: time.Now()},
		{Name: "Shape", CreatedAt: time.Now()},
		{Name: "Color", CreatedAt: time.Now()},
		{Name: "Origin", CreatedAt: time.Now()},
		{Name: "Material", CreatedAt: time.Now()},
		{Name: "Purpose", CreatedAt: time.Now()},
	}

	for _, c := range categoryLabels {
		label := c
		label.Slug = utils.SlugifyText(label.Name)

		if err := db.FirstOrCreate(&label, models.CategoryLabel{Name: label.Name}).Error; err != nil {
			log.Printf("Gagal seeding category_label %s: %v\n", label.Name, err)
		} else {
			if err := db.Model(&label).Update("slug", label.Slug).Error; err != nil {
				log.Printf("Gagal update slug category_label %s: %v\n", label.Name, err)
			}
		}
	}

	log.Println("[SEEDER] ⚙️  category labels seeded")
}
