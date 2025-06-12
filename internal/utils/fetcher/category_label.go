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

func GetAllCategoryLabels(ctx context.Context) ([]models.CategoryLabel, error) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal ambil Redis client:", err)
		rdb = nil
	}

	if rdb != nil {
		data, err := rdb.HGetAll(ctx, "category_label:id").Result()
		if err == nil && len(data) > 0 {
			var labels []models.CategoryLabel
			for idStr, jsonStr := range data {
				var label models.CategoryLabel
				if err := json.Unmarshal([]byte(jsonStr), &label); err != nil {
					log.Printf("[CACHE] ⚠️ Gagal decode JSON ID %s: %v", idStr, err)
					continue
				}
				if idInt, convErr := strconv.Atoi(idStr); convErr == nil {
					label.ID = uint(idInt)
				}
				labels = append(labels, label)
			}
			log.Printf("[CACHE] ✅ Berhasil ambil %d category label dari Redis", len(labels))
			return labels, nil
		}
		log.Println("[CACHE] ℹ️ Redis kosong atau error, fallback ke DB")
	}

	var labels []models.CategoryLabel
	if err := database.DB.Order("name ASC").Find(&labels).Error; err != nil {
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d category label dari database", len(labels))

	if rdb != nil {
		cacheMap := make(map[string]string)
		for _, label := range labels {
			jsonStr, err := json.Marshal(label)
			if err != nil {
				log.Printf("[CACHE] ⚠️ Gagal encode category label ID %d: %v", label.ID, err)
				continue
			}
			cacheMap[strconv.Itoa(int(label.ID))] = string(jsonStr)
		}
		if err := rdb.HSet(ctx, "category_label:id", cacheMap).Err(); err != nil {
			log.Println("[CACHE] ⚠️ Gagal menyimpan category label ke Redis:", err)
		} else {
			log.Printf("[CACHE] ✅ Category label disimpan ke Redis (%d item)", len(cacheMap))
		}
	}

	return labels, nil
}
