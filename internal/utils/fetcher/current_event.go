package fetcher

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func GetCurrentEvent(ctx context.Context) (*models.CurrentEvent, error) {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	cacheKey := "current_event"

	if err == nil {
		jsonStr, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil && jsonStr != "" {
			var ce models.CurrentEvent
			if err := json.Unmarshal([]byte(jsonStr), &ce); err == nil {
				log.Println("[CACHE] ✅ CurrentEvent ditemukan di Redis")
				return &ce, nil
			}
			log.Println("[CACHE] ⚠️ Gagal decode CurrentEvent dari Redis:", err)
		}
	}

	var currentEvent models.CurrentEvent
	now := time.Now()

	err = database.DB.
		Preload("EventType").
		Where("start <= ?", now).
		Order("created_at DESC").
		First(&currentEvent).Error

	if err != nil {
		log.Println("[DB] ❌ Tidak ada CurrentEvent aktif:", err)
		return nil, err
	}

	log.Println("[DB] ✅ CurrentEvent aktif ditemukan di DB:", currentEvent.ID)

	if rdb != nil {
		jsonVal, err := json.Marshal(currentEvent)
		if err == nil {
			if err := rdb.Set(ctx, cacheKey, jsonVal, 24*time.Hour).Err(); err == nil {
				log.Println("[CACHE] ✅ CurrentEvent disimpan ke Redis")
			}
		}
	}

	return &currentEvent, nil
}
