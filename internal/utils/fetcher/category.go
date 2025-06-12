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

func GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal ambil Redis client:", err)
		rdb = nil
	}

	if rdb != nil {
		data, err := rdb.HGetAll(ctx, "category:id").Result()
		if err == nil && len(data) > 0 {
			var categories []models.Category
			for idStr, jsonStr := range data {
				var cat models.Category
				if err := json.Unmarshal([]byte(jsonStr), &cat); err != nil {
					log.Printf("[CACHE] ⚠️ Gagal decode JSON ID %s: %v", idStr, err)
					continue
				}
				if idInt, convErr := strconv.Atoi(idStr); convErr == nil {
					cat.ID = uint(idInt)
				}
				categories = append(categories, cat)
			}
			log.Printf("[CACHE] ✅ Berhasil ambil %d category dari Redis", len(categories))
			return categories, nil
		}
		log.Println("[CACHE] ℹ️ Redis kosong atau error, fallback ke DB")
	}

	var categories []models.Category
	if err := database.DB.Preload("CategoryLabel").Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d category dari database", len(categories))

	if rdb != nil {
		cacheMap := make(map[string]string)
		for _, cat := range categories {
			jsonStr, err := json.Marshal(cat)
			if err != nil {
				log.Printf("[CACHE] ⚠️ Gagal encode category ID %d: %v", cat.ID, err)
				continue
			}
			cacheMap[strconv.Itoa(int(cat.ID))] = string(jsonStr)
		}
		if err := rdb.HSet(ctx, "category:id", cacheMap).Err(); err != nil {
			log.Println("[CACHE] ⚠️ Gagal menyimpan category ke Redis:", err)
		} else {
			log.Printf("[CACHE] ✅ Category disimpan ke Redis (%d item)", len(cacheMap))
		}
	}

	return categories, nil
}
