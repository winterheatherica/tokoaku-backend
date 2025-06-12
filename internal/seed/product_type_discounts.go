package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedProductTypeDiscounts(db *gorm.DB) {
	productTypeIDs := []uint{3, 23, 24}
	discountIDs := []uint{1, 2, 3, 4, 5}

	var data []models.ProductTypeDiscount
	now := time.Now()

	for _, ptID := range productTypeIDs {
		for _, dID := range discountIDs {
			data = append(data, models.ProductTypeDiscount{
				ProductTypeID: ptID,
				DiscountID:    dID,
				CreatedAt:     now,
			})
		}
	}

	for _, d := range data {
		if err := db.FirstOrCreate(&d, models.ProductTypeDiscount{
			ProductTypeID: d.ProductTypeID,
			DiscountID:    d.DiscountID,
		}).Error; err != nil {
			log.Printf("Gagal seeding ProductTypeDiscount: ProductTypeID=%d, DiscountID=%d, error=%v\n", d.ProductTypeID, d.DiscountID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  ProductTypeDiscounts seeded")
}
