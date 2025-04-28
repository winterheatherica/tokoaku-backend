package cloudinary

import (
	"fmt"
	"log"
	"os"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

type CloudinaryDynamic struct {
	CloudName string
	APIKey    string
	APISecret string
}

var cloudinaryConfigs map[string]*CloudinaryDynamic

func LoadAllCloudinaryConfigs(db *gorm.DB) error {
	var services []models.CloudService

	err := db.Where("provider_id = ?", 11).Find(&services).Error
	if err != nil {
		return fmt.Errorf("failed to load cloud services: %w", err)
	}

	cloudinaryConfigs = make(map[string]*CloudinaryDynamic)

	for _, service := range services {

		cloudName := os.Getenv(fmt.Sprintf("%s_CLOUDINARY_CLOUD_NAME", service.EnvKeyPrefix))
		apiKey := os.Getenv(fmt.Sprintf("%s_CLOUDINARY_API_KEY", service.EnvKeyPrefix))
		apiSecret := os.Getenv(fmt.Sprintf("%s_CLOUDINARY_API_SECRET", service.EnvKeyPrefix))

		if cloudName == "" || apiKey == "" || apiSecret == "" {
			log.Printf("Missing environment variable untuk prefix: %s, skip...", service.EnvKeyPrefix)
			continue
		}

		cloudinaryConfigs[service.EnvKeyPrefix] = &CloudinaryDynamic{
			CloudName: cloudName,
			APIKey:    apiKey,
			APISecret: apiSecret,
		}

		log.Printf("[SERVICE]: ⚙️  Cloudinary prefix: %s connected", service.EnvKeyPrefix)
	}

	return nil
}

func GetCloudinaryConfig(prefix string) (*CloudinaryDynamic, error) {
	cfg, ok := cloudinaryConfigs[prefix]
	if !ok {
		return nil, fmt.Errorf("cloudinary config for prefix %s not found", prefix)
	}
	return cfg, nil
}
