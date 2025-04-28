package volatile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func refreshCategories() {
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

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		log.Println("[CACHE] 🔄 Refresh Categories (Key ID, Value Name)")

		var categories []models.Category
		if err := database.DB.Order("id ASC").Find(&categories).Error; err != nil {
			log.Println("[CACHE] Gagal refresh Categories:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, cat := range categories {
			key := fmt.Sprintf("category:%d", cat.ID)
			value := cat.Name

			if err := redisClient.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
				log.Println("[CACHE] Gagal set Category:", err)
			}
		}

		<-ticker.C
	}
}
