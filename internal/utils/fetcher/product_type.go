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

func GetAllProductTypes(ctx context.Context) ([]models.ProductType, error) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal ambil Redis client:", err)
		rdb = nil
	}

	if rdb != nil {
		data, err := rdb.HGetAll(ctx, "product_type:id").Result()
		if err == nil && len(data) > 0 {
			var productTypes []models.ProductType
			for idStr, jsonStr := range data {
				var pt models.ProductType
				if err := json.Unmarshal([]byte(jsonStr), &pt); err != nil {
					log.Printf("[CACHE] ⚠️ Gagal decode JSON ID %s: %v", idStr, err)
					continue
				}
				if idInt, convErr := strconv.Atoi(idStr); convErr == nil {
					pt.ID = uint(idInt)
				}
				productTypes = append(productTypes, pt)
			}
			log.Printf("[CACHE] ✅ Berhasil ambil %d product type dari Redis", len(productTypes))
			return productTypes, nil
		}
		log.Println("[CACHE] ℹ️ Redis kosong atau error, fallback ke DB")
	}

	var types []models.ProductType
	if err := database.DB.Order("name ASC").Find(&types).Error; err != nil {
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d product type dari database", len(types))

	if rdb != nil {
		cacheMap := make(map[string]string)
		for _, pt := range types {
			jsonStr, err := json.Marshal(pt)
			if err != nil {
				log.Printf("[CACHE] ⚠️ Gagal encode product type ID %d: %v", pt.ID, err)
				continue
			}
			cacheMap[strconv.Itoa(int(pt.ID))] = string(jsonStr)
		}
		if err := rdb.HSet(ctx, "product_type:id", cacheMap).Err(); err != nil {
			log.Println("[CACHE] ⚠️ Gagal menyimpan product type ke Redis:", err)
		} else {
			log.Printf("[CACHE] ✅ Product type disimpan ke Redis (%d item)", len(cacheMap))
		}
	}

	return types, nil
}
