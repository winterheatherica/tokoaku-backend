package volatile

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/config"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetAllVolatileRedisClients(ctx context.Context) ([]*redis.Client, error) {
	var services []models.CloudService
	if err := database.DB.Where("usage_for = ?", "Volatile Cache").Order("storage_usage ASC").Find(&services).Error; err != nil {
		return nil, err
	}

	var clients []*redis.Client
	for _, svc := range services {
		rdb, err := config.GetRedisClient(svc.EnvKeyPrefix)
		if err != nil {
			log.Printf("[CACHE] ⚠️ Gagal ambil Redis client untuk prefix %s: %v", svc.EnvKeyPrefix, err)
			continue
		}
		clients = append(clients, rdb)
	}
	return clients, nil
}
