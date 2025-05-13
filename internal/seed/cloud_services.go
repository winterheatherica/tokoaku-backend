package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedCloudServices(db *gorm.DB) {
	cloudServices := []models.CloudService{
		{
			Name:          "Erika Cloudinary Storage",
			ProviderID:    3,
			EnvKeyPrefix:  "ERIKA",
			UsageFor:      "Private Image",
			StorageUsage:  0,
			LastCheckedAt: time.Now(),
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Keqing Cloudinary Storage",
			ProviderID:    3,
			EnvKeyPrefix:  "KEQING",
			UsageFor:      "Public Image",
			StorageUsage:  0,
			LastCheckedAt: time.Now(),
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Nilou Upstash Redis Cache",
			ProviderID:    4,
			EnvKeyPrefix:  "NILOU",
			UsageFor:      "Persistent Cache",
			StorageUsage:  0,
			LastCheckedAt: time.Now(),
			CreatedAt:     time.Now(),
		},
		{
			Name:          "Hu Tao Upstash Redis Cache",
			ProviderID:    4,
			EnvKeyPrefix:  "HUTAO",
			UsageFor:      "Volatile Cache",
			StorageUsage:  0,
			LastCheckedAt: time.Now(),
			CreatedAt:     time.Now(),
		},
	}

	for _, cs := range cloudServices {
		if err := db.FirstOrCreate(&cs, models.CloudService{EnvKeyPrefix: cs.EnvKeyPrefix}).Error; err != nil {
			log.Printf("Gagal seeding cloud_service dengan prefix %s: %v\n", cs.EnvKeyPrefix, err)
		}
	}

	log.Println("[SEEDER] ⚙️  cloud services seeded")
}
