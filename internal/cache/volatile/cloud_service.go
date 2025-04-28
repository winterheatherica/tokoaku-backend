package volatile

import (
	"context"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func refreshCloudServices() {
	ctx := context.Background()

	prefix, err := utils.GetVolatileRedisPrefix()
	if err != nil {
		log.Println("[CACHE] Gagal ambil volatile prefix:", err)
		return
	}

	redisClient, err := redis.GetRedisClient(prefix)
	if err != nil {
		log.Println("[CACHE] Gagal ambil Redis client:", err)
		return
	}

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		log.Println("[CACHE] 🔄 Refresh CloudServices (SET Mode)")

		var services []models.CloudService
		if err := database.DB.Order("env_key_prefix ASC").Find(&services).Error; err != nil {
			log.Println("[CACHE] Gagal refresh CloudServices:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, svc := range services {
			key := formatUsageForKey(svc.UsageFor)
			value := svc.EnvKeyPrefix

			if err := redisClient.SAdd(ctx, key, value).Err(); err != nil {
				log.Println("[CACHE] Gagal SADD CloudService:", err)
			}
		}

		<-ticker.C
	}
}

func formatUsageForKey(usage string) string {
	return utils.ToSnakeCase(usage)
}
