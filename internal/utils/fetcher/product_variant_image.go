package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

type VariantImageLite struct {
	IsVariantCover bool `json:"is_variant_cover"`
}

func GetVariantCoverImage(ctx context.Context, variantID string) (*models.ProductVariantImage, error) {
	cacheKey := fmt.Sprintf("variant:images:%s", variantID)

	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err == nil {
		dataMap, err := rdb.HGetAll(ctx, cacheKey).Result()
		if err == nil && len(dataMap) > 0 {
			for url, val := range dataMap {
				var lite VariantImageLite
				if err := json.Unmarshal([]byte(val), &lite); err == nil && lite.IsVariantCover {
					log.Printf("[CACHE] ✅ Cover variant %s ditemukan di Redis", variantID)
					return &models.ProductVariantImage{
						ProductVariantID: variantID,
						ImageURL:         url,
						IsVariantCover:   true,
					}, nil
				}
			}
		}
	}

	var image models.ProductVariantImage
	if err := database.DB.
		Where("product_variant_id = ? AND is_variant_cover = true", variantID).
		First(&image).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil cover variant %s dari DB: %v", variantID, err)
		return nil, err
	}
	log.Printf("[DB] ✅ Cover variant %s ditemukan di DB", variantID)

	if rdb != nil {
		_, _ = GetAllVariantImages(ctx, variantID)
	}

	return &image, nil
}

func GetAllVariantImages(ctx context.Context, variantID string) ([]models.ProductVariantImage, error) {
	cacheKey := fmt.Sprintf("variant:images:%s", variantID)

	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err == nil {
		dataMap, err := rdb.HGetAll(ctx, cacheKey).Result()
		if err == nil && len(dataMap) > 0 {
			var images []models.ProductVariantImage
			for url, val := range dataMap {
				var lite VariantImageLite
				if err := json.Unmarshal([]byte(val), &lite); err == nil {
					images = append(images, models.ProductVariantImage{
						ProductVariantID: variantID,
						ImageURL:         url,
						IsVariantCover:   lite.IsVariantCover,
					})
				}
			}
			log.Printf("[CACHE] ✅ Semua image variant %s ditemukan di Redis", variantID)
			return images, nil
		}
	}

	var images []models.ProductVariantImage
	if err := database.DB.
		Where("product_variant_id = ?", variantID).
		Order("created_at ASC").
		Find(&images).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil semua image variant %s: %v", variantID, err)
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d image variant %s dari DB", len(images), variantID)

	if rdb != nil && len(images) > 0 {
		entries := make(map[string]string)
		for _, img := range images {
			val, _ := json.Marshal(VariantImageLite{
				IsVariantCover: img.IsVariantCover,
			})
			entries[img.ImageURL] = string(val)
		}

		_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, cacheKey, entries)
			pipe.Expire(ctx, cacheKey, 3*time.Hour)
			return nil
		})
		if err != nil {
			log.Printf("[CACHE] ⚠️ Gagal simpan variant %s ke Redis (HSET): %v", variantID, err)
		} else {
			log.Printf("[CACHE] ✅ Semua image variant %s disimpan ke Redis (HSET, TTL 3 jam)", variantID)
		}
	}

	return images, nil
}
