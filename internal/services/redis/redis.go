package redis

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/config"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func LoadAllRedisConfigs(db *gorm.DB) error {
	var services []models.CloudService
	if err := db.Preload("Provider").Find(&services).Error; err != nil {
		return fmt.Errorf("failed to load Redis cloud services: %w", err)
	}

	for _, service := range services {
		if service.Provider.Name != "Upstash (Redis)" {
			continue
		}

		prefix := service.EnvKeyPrefix
		envKey := fmt.Sprintf("%s_REDIS_URL", prefix)
		if os.Getenv(envKey) == "" {
			log.Printf("⚠️ Missing REDIS_URL for prefix: %s", prefix)
			continue
		}

		if err := config.LoadRedisClient(prefix); err != nil {
			log.Printf("⚠️ Failed to load Redis client for %s: %v", prefix, err)
			continue
		}
	}

	return nil
}

func GetRedisClient(prefix string) (*redis.Client, error) {
	return config.GetRedisClient(prefix)
}
