package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedShippingOptions(db *gorm.DB) {
	shippingOptions := []models.ShippingOption{
		{
			ID:                 1,
			CourierName:        "JNE",
			CourierServiceName: "REG",
			Fee:                18000,
			EstimatedTime:      "2-3 Hari",
			ServiceType:        "Regular",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 2,
			CourierName:        "JNE",
			CourierServiceName: "YES",
			Fee:                35000,
			EstimatedTime:      "1 Hari",
			ServiceType:        "Express",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 3,
			CourierName:        "SiCepat",
			CourierServiceName: "BEST",
			Fee:                29000,
			EstimatedTime:      "1 Hari",
			ServiceType:        "Express",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 4,
			CourierName:        "SiCepat",
			CourierServiceName: "REG",
			Fee:                17000,
			EstimatedTime:      "2-3 Hari",
			ServiceType:        "Regular",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 5,
			CourierName:        "Anteraja",
			CourierServiceName: "Regular",
			Fee:                15000,
			EstimatedTime:      "2-4 Hari",
			ServiceType:        "Regular",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 6,
			CourierName:        "Anteraja",
			CourierServiceName: "Next Day",
			Fee:                25000,
			EstimatedTime:      "1 Hari",
			ServiceType:        "Express",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 7,
			CourierName:        "J&T Express",
			CourierServiceName: "EZ",
			Fee:                16000,
			EstimatedTime:      "2-3 Hari",
			ServiceType:        "Regular",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 100,
			CourierName:        "None",
			CourierServiceName: "None",
			Fee:                0,
			EstimatedTime:      "None",
			ServiceType:        "None",
			CreatedAt:          time.Now(),
		},
	}

	for _, s := range shippingOptions {
		if err := db.FirstOrCreate(&s, models.ShippingOption{ID: s.ID}).Error; err != nil {
			log.Printf("Gagal seeding shipping_option ID %d: %v\n", s.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  shipping options seeded")
}
