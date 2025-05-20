package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedPromos(db *gorm.DB) {
	promos := []models.Promo{
		{
			Name:          "Test",
			Code:          "TEST2025",
			Description:   "Promo test dari tanggal 28 April 2025 sampai 12 Desember 2025",
			ValueTypeID:   1,
			Value:         10,
			MinPriceValue: 100000,
			MaxValue:      50000,
			StartAt:       time.Date(2025, 4, 28, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 12, 12, 23, 59, 59, 0, time.UTC),
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Mid Year Sale",
			Code:          "MIDYEAR25",
			Description:   "Diskon 25% untuk semua produk selama pertengahan tahun.",
			ValueTypeID:   1,
			Value:         25,
			MinPriceValue: 50000,
			MaxValue:      75000,
			StartAt:       time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Flash Sale 7.7",
			Code:          "FLASH77",
			Description:   "Promo spesial 7 Juli! Dapatkan potongan langsung Rp30.000.",
			ValueTypeID:   2,
			Value:         30000,
			MinPriceValue: 0,
			MaxValue:      0,
			StartAt:       time.Date(2025, 7, 7, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 7, 7, 23, 59, 59, 0, time.UTC),
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Independence Day Sale",
			Code:          "MERDEKA45",
			Description:   "Diskon 45% untuk merayakan Hari Kemerdekaan.",
			ValueTypeID:   1,
			Value:         45,
			MinPriceValue: 100000,
			MaxValue:      100000,
			StartAt:       time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 8, 17, 23, 59, 59, 0, time.UTC),
			CreatedAt:     time.Now(),
		},
	}

	for _, p := range promos {
		if err := db.FirstOrCreate(&p, models.Promo{Code: p.Code}).Error; err != nil {
			log.Printf("Gagal seeding promo kode %s: %v\n", p.Code, err)
		}
	}

	log.Println("[SEEDER] ⚙️  promos seeded")
}
