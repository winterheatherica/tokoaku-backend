package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:slug:%s", slug)

	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err == nil {
		if jsonStr, err := rdb.Get(ctx, cacheKey).Result(); err == nil && jsonStr != "" {
			var cached models.Product
			if err := json.Unmarshal([]byte(jsonStr), &cached); err == nil {
				log.Printf("[CACHE] ✅ Product %s ditemukan di Redis", slug)
				return &cached, nil
			}
			log.Printf("[CACHE] ⚠️ Gagal decode cache product %s: %v", slug, err)
		}
	}

	var product models.Product
	if err := database.DB.
		WithContext(ctx).
		Preload("ProductType").
		Preload("ProductForm").
		Preload("Variants").
		Where("slug = ?", slug).
		First(&product).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil product %s: %v", slug, err)
		return nil, err
	}

	log.Printf("[DB] ✅ Product %s berhasil diambil dari DB", slug)
	log.Printf("[DEBUG] ✅ Product %s punya %d variant(s)", product.ID, len(product.Variants))

	if rdb != nil {
		if jsonData, err := json.Marshal(product); err == nil {
			if err := rdb.Set(ctx, cacheKey, jsonData, 1*time.Hour).Err(); err == nil {
				log.Printf("[CACHE] ✅ Product %s disimpan ke Redis", slug)
			}
		}
	}

	return &product, nil
}
