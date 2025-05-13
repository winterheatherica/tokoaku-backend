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

func storeProductTypeToRedis(ctx context.Context, rdb *redis.Client, pt models.ProductType) {
	idKey := "product_type:id"
	nameKey := "product_type:name"

	value := map[string]interface{}{
		"name":          pt.Name,
		"slug":          pt.Slug,
		"value_type_id": pt.ValueTypeID,
		"value":         pt.Value,
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal encode JSON ProductType ID %d: %v", pt.ID, err)
		return
	}

	if err := rdb.HSet(ctx, idKey, fmt.Sprintf("%d", pt.ID), jsonValue).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s → ID %d: %v", idKey, pt.ID, err)
	}

	if err := rdb.HSet(ctx, nameKey, pt.Name, pt.ID).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s → Name %s: %v", nameKey, pt.Name, err)
	}
}

func refreshProductTypes() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshProductTypes (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval24h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh ProductTypes (Persistent)")

			var productTypes []models.ProductType
			err := database.DB.Order("id ASC").Find(&productTypes).Error
			if err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data ProductTypes:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, pt := range productTypes {
				storeProductTypeToRedis(ctx, rdb, pt)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d product type ke Redis (Persistent)", len(productTypes))
			<-ticker.C
		}
	}()
}
