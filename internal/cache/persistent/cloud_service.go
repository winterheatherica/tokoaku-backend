package persistent

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func storeCloudServiceToRedis(ctx context.Context, rdb *redis.Client, svc models.CloudService) {
	key := utils.ToSnakeCase(svc.UsageFor)
	value := map[string]interface{}{
		"id":             svc.ID,
		"env_key_prefix": svc.EnvKeyPrefix,
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal encode CloudService ID %d: %v", svc.ID, err)
		return
	}

	if err := rdb.SAdd(ctx, key, jsonValue).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal SADD ke Redis untuk key %s: %v", key, err)
	}
}

func refreshCloudServices() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshCloudServices (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval1h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh CloudServices (Persistent)")

			var services []models.CloudService
			err := database.DB.Order("env_key_prefix ASC").Find(&services).Error
			if err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data CloudServices:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, svc := range services {
				storeCloudServiceToRedis(ctx, rdb, svc)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d cloud service ke Redis (Persistent)", len(services))
			<-ticker.C
		}
	}()
}
