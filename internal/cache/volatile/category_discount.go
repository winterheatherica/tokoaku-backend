package volatile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func storeCategoryDiscountToRedis(ctx context.Context, rdb *redis.Client, cd models.CategoryDiscount) {
	key := fmt.Sprintf("category_discount:%s", cd.Category.Slug)
	value := cd.Discount.Value

	if err := rdb.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal set CategoryDiscount %s: %v", key, err)
	}
}

func refreshCategoryDiscounts() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshCategoryDiscounts (Volatile)")

		ctx := context.Background()

		rdb, err := volatile.GetVolatileRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval6h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh CategoryDiscounts (Volatile)")

			var categoryDiscounts []models.CategoryDiscount
			err := database.DB.
				Preload("Category").
				Preload("Discount").
				Order("category_id ASC, discount_id ASC").
				Find(&categoryDiscounts).Error

			if err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data CategoryDiscounts:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, cd := range categoryDiscounts {
				storeCategoryDiscountToRedis(ctx, rdb, cd)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d data diskon kategori ke Redis (Volatile)", len(categoryDiscounts))
			<-ticker.C
		}
	}()
}
