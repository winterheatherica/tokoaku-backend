package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedDiscounts(db *gorm.DB) {
	discounts := []models.Discount{
		{
			Name:        "Diskon Ramadhan 2025",
			ValueTypeID: 1,
			Value:       20,
			SponsorID:   2,
			Start:       time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
			End:         time.Date(2025, 4, 10, 23, 59, 59, 0, time.UTC),
			Slug:        "diskon-ramadhan-2025",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "Diskon Bulan April 2025",
			ValueTypeID: 1,
			Value:       15,
			SponsorID:   2,
			Start:       time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			End:         time.Date(2025, 4, 30, 23, 59, 59, 0, time.UTC),
			Slug:        "diskon-bulan-april-2025",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "Diskon Bulan Mei 2025",
			ValueTypeID: 2,
			Value:       50000,
			SponsorID:   2,
			Start:       time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			End:         time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC),
			Slug:        "diskon-bulan-mei-2025",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "Diskon Natal 2025",
			ValueTypeID: 1,
			Value:       25,
			SponsorID:   2,
			Start:       time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			End:         time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			Slug:        "diskon-natal-2025",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "Diskon Go-Toubun no Hanayome 2025",
			ValueTypeID: 2,
			Value:       55000,
			SponsorID:   2,
			Start:       time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			End:         time.Date(2025, 5, 6, 23, 59, 59, 0, time.UTC),
			Slug:        "diskon-go-toubun-no-hanayome-2025",
			CreatedAt:   time.Now(),
		},
	}

	for _, d := range discounts {
		if err := db.FirstOrCreate(&d, models.Discount{Slug: d.Slug}).Error; err != nil {
			log.Printf("Gagal seeding discount slug %s: %v\n", d.Slug, err)
		}
	}

	log.Println("[SEEDER] ⚙️  discounts seeded")
}
