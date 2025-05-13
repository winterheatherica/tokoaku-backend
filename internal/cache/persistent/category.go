package persistent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func storeCategoryToRedis(ctx context.Context, rdb *redis.Client, cat models.Category) {
	idKey := "category:id"
	nameKey := "category:name"

	value := map[string]string{
		"name": cat.Name,
		"slug": cat.Slug,
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal Marshal JSON Category ID %d: %v", cat.ID, err)
		return
	}

	if err := rdb.HSet(ctx, idKey, fmt.Sprintf("%d", cat.ID), jsonValue).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s -> ID %d: %v", idKey, cat.ID, err)
	}

	if err := rdb.HSet(ctx, nameKey, cat.Name, cat.ID).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s -> Name %s: %v", nameKey, cat.Name, err)
	}
}

func refreshCategories() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshCategories (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval24h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh Categories (Persistent)")

			var categories []models.Category
			if err := database.DB.Order("id ASC").Find(&categories).Error; err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data kategori dari database:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, cat := range categories {
				storeCategoryToRedis(ctx, rdb, cat)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d kategori ke Redis (Persistent)", len(categories))
			<-ticker.C
		}
	}()
}
