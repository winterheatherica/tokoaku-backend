package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	"gorm.io/gorm"
)

func SeedDiscounts(db *gorm.DB) {
	discounts := []models.Discount{
		{
			Name:          "Ramadhan Sale 2025",
			Description:   "Get up to 20% off on selected items. Limited time offer!",
			ValueTypeID:   1,
			Value:         20,
			SponsorID:     2,
			StartAt:       time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 4, 10, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
		{
			Name:          "April Fools 2025",
			Description:   "Get up to 15% off on selected items. Limited time offer!",
			ValueTypeID:   1,
			Value:         15,
			SponsorID:     2,
			StartAt:       time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 4, 30, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
		{
			Name:          "We Are From 27th May",
			Description:   "Get up to Rp 2000 off on selected items. Limited time offer!",
			ValueTypeID:   2,
			Value:         2000,
			SponsorID:     2,
			StartAt:       time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Christmas Sale 2025",
			Description:   "Get up to 25% off on selected items. Limited time offer!",
			ValueTypeID:   1,
			Value:         25,
			SponsorID:     2,
			StartAt:       time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Go-Toubun no Hanayome 2025",
			Description:   "Get up to Rp 3000 off on selected items. Limited time offer!",
			ValueTypeID:   2,
			Value:         3000,
			SponsorID:     2,
			StartAt:       time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 5, 6, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
		{
			Name:          "June Sale 2025",
			Description:   "Get up to 2% off on selected items. Limited time offer!",
			ValueTypeID:   1,
			Value:         2,
			SponsorID:     2,
			StartAt:       time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
		{
			Name:          "May to June",
			Description:   "Get up to 3% off on selected items. Limited time offer!",
			ValueTypeID:   1,
			Value:         3,
			SponsorID:     2,
			StartAt:       time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			EndAt:         time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
			ImageCoverURL: "https://asset.cloudinary.com/dokzc5ogk/cd40085bc7c20b476e24fcbcf9c83881",
			CloudImageID:  2,
			CreatedAt:     time.Now(),
		},
	}

	log.Printf("[DEBUG] Ditemukan %d total diskon dari event aktif", len(discounts))

	for i := range discounts {
		if discounts[i].Slug == "" {
			discounts[i].Slug = utils.SlugifyText(discounts[i].Name)
		}

		if err := db.FirstOrCreate(&discounts[i], models.Discount{Slug: discounts[i].Slug}).Error; err != nil {
			log.Printf("Gagal seeding discount slug %s: %v\n", discounts[i].Slug, err)
		}
	}

	log.Println("[SEEDER] ⚙️  discounts seeded")
}
