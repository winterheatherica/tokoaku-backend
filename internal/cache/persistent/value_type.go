package persistent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func storeValueTypeToRedis(ctx context.Context, rdb *redis.Client, vt models.ValueType) {
	idKey := "value_type:id"
	nameKey := "value_type:name"

	if err := rdb.HSet(ctx, idKey, fmt.Sprintf("%d", vt.ID), vt.Name).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s → ID %d: %v", idKey, vt.ID, err)
	}

	if err := rdb.HSet(ctx, nameKey, vt.Name, vt.ID).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s → Name %s: %v", nameKey, vt.Name, err)
	}
}

func refreshValueTypes() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshValueTypes (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval24h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh ValueTypes (Persistent)")

			var valueTypes []models.ValueType
			err := database.DB.Order("id ASC").Find(&valueTypes).Error
			if err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data ValueTypes:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, vt := range valueTypes {
				storeValueTypeToRedis(ctx, rdb, vt)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d value type ke Redis (Persistent)", len(valueTypes))
			<-ticker.C
		}
	}()
}
