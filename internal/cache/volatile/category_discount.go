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

func refreshCategoryDiscounts() {
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
		log.Println("[CACHE] 🔄 Refresh CategoryDiscounts")

		var categoryDiscounts []models.CategoryDiscount
		if err := database.DB.
			Preload("Category").
			Preload("Discount").
			Order("category_id ASC, discount_id ASC").
			Find(&categoryDiscounts).Error; err != nil {
			log.Println("[CACHE] Gagal refresh CategoryDiscounts:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, cd := range categoryDiscounts {
			key := fmt.Sprintf("category_discount:%s", cd.Category.Slug)
			value := cd.Discount.Value

			if err := redisClient.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
				log.Println("[CACHE] Gagal set CategoryDiscount:", err)
			}
		}

		<-ticker.C
	}
}
