package fetcher

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func GetAllEventTypes(ctx context.Context) ([]models.EventType, error) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal ambil Redis client:", err)
		rdb = nil
	}

	const cacheKey = "event_type:id"

	if rdb != nil {
		data, err := rdb.HGetAll(ctx, cacheKey).Result()
		if err == nil && len(data) > 0 {
			var result []models.EventType
			for idStr, jsonStr := range data {
				var et models.EventType
				if err := json.Unmarshal([]byte(jsonStr), &et); err != nil {
					log.Printf("[CACHE] ⚠️ Gagal decode JSON ID %s: %v", idStr, err)
					continue
				}
				if idInt, err := strconv.Atoi(idStr); err == nil {
					et.ID = uint(idInt)
				}
				result = append(result, et)
			}
			log.Printf("[CACHE] ✅ Berhasil ambil %d event type dari Redis", len(result))
			return result, nil
		}
		log.Println("[CACHE] ℹ️ Redis kosong atau error, fallback ke DB")
	}

	var eventTypes []models.EventType
	if err := database.DB.Order("id ASC").Find(&eventTypes).Error; err != nil {
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d event type dari database", len(eventTypes))

	if rdb != nil {
		entries := make(map[string]string)
		for _, et := range eventTypes {
			jsonVal, err := json.Marshal(et)
			if err != nil {
				continue
			}
			entries[strconv.Itoa(int(et.ID))] = string(jsonVal)
		}
		if err := rdb.HSet(ctx, cacheKey, entries).Err(); err != nil {
			log.Println("[CACHE] ⚠️ Gagal simpan event type ke Redis:", err)
		} else {
			log.Printf("[CACHE] ✅ Event type disimpan ke Redis (%d item)", len(entries))
		}
	}

	return eventTypes, nil
}
