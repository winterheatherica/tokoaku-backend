package cloudinary

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

var (
	cloudinaryConfigs = make(map[string]*CloudinaryConfig)
	configLock        sync.RWMutex
)

func LoadAllCloudinaryConfigs(db *gorm.DB) error {
	var services []models.CloudService

	if err := db.Preload("Provider").Find(&services).Error; err != nil {
		return fmt.Errorf("failed to load cloud services: %w", err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	for _, service := range services {
		if service.Provider.Name != "Cloudinary" {
			continue
		}

		cfg := loadFromEnv(service.EnvKeyPrefix)
		if cfg == nil {
			log.Printf("⚠️ Missing ENV for prefix: %s, skipped", service.EnvKeyPrefix)
			continue
		}

		cloudinaryConfigs[service.EnvKeyPrefix] = cfg
		log.Printf("[CLOUDINARY] ✅ Loaded config for %s", service.EnvKeyPrefix)
	}
	return nil
}

func loadFromEnv(prefix string) *CloudinaryConfig {
	cloudName := os.Getenv(fmt.Sprintf("%s_CLOUDINARY_CLOUD_NAME", prefix))
	apiKey := os.Getenv(fmt.Sprintf("%s_CLOUDINARY_API_KEY", prefix))
	apiSecret := os.Getenv(fmt.Sprintf("%s_CLOUDINARY_API_SECRET", prefix))

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil
	}

	return &CloudinaryConfig{
		CloudName: cloudName,
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

func GetCloudinaryConfig(prefix string) (*CloudinaryConfig, error) {
	configLock.RLock()
	defer configLock.RUnlock()

	cfg, ok := cloudinaryConfigs[prefix]
	if !ok {
		return nil, fmt.Errorf("cloudinary config for prefix %s not found", prefix)
	}
	return cfg, nil
}
