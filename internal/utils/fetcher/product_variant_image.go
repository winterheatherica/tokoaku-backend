package fetcher

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

type VariantImageLite struct {
	ImageURL       string `json:"image_url"`
	CloudImageID   uint   `json:"cloud_image_id"`
	IsVariantCover bool   `json:"is_variant_cover"`
}

func GetVariantCoverImage(ctx context.Context, productID, variantID string) (*models.ProductVariantImage, error) {
	images, err := getVariantImagesFromRedis(ctx, variantID)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal ambil cover dari Redis: %v", err)
		return nil, err
	}

	if len(images) == 0 {
		log.Printf("[CACHE] ⚠️ Tidak ada cover image variant %s di Redis", variantID)
		return nil, nil
	}

	for _, img := range images {
		if img.IsVariantCover {
			log.Printf("[CACHE] ✅ Cover variant %s ditemukan di Redis", variantID)
			return &models.ProductVariantImage{
				ProductVariantID: variantID,
				ImageURL:         img.ImageURL,
				CloudImageID:     img.CloudImageID,
				IsVariantCover:   img.IsVariantCover,
			}, nil
		}
	}

	log.Printf("[CACHE] ⚠️ Redis valid, tapi tidak ditemukan cover image untuk variant %s", variantID)
	return nil, nil
}

func GetAllVariantImages(ctx context.Context, variantID string) ([]models.ProductVariantImage, error) {
	images, err := getVariantImagesFromRedis(ctx, variantID)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal ambil gambar variant %s dari Redis: %v", variantID, err)
		return nil, err
	}

	if images == nil {
		log.Printf("[CACHE] ❓ Belum ada cache untuk variant %s → Query ke DB", variantID)

		var dbImages []models.ProductVariantImage
		if err := database.DB.
			WithContext(ctx).
			Where("product_variant_id = ?", variantID).
			Order("created_at ASC").
			Find(&dbImages).Error; err != nil {
			log.Printf("[DB] ❌ Gagal ambil gambar dari DB: %v", err)
			return nil, err
		}

		_ = cacheImagesToRedis(ctx, variantID, dbImages)
		return dbImages, nil
	}

	if len(images) == 0 {
		return []models.ProductVariantImage{}, nil
	}

	var result []models.ProductVariantImage
	for _, img := range images {
		result = append(result, models.ProductVariantImage{
			ProductVariantID: variantID,
			ImageURL:         img.ImageURL,
			CloudImageID:     img.CloudImageID,
			IsVariantCover:   img.IsVariantCover,
		})
	}
	return result, nil
}

func CacheVariantImageFromDB(ctx context.Context, variantID string) error {
	var images []models.ProductVariantImage
	if err := database.DB.
		WithContext(ctx).
		Where("product_variant_id = ?", variantID).
		Order("created_at ASC").
		Find(&images).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil gambar variant %s untuk cache: %v", variantID, err)
		return err
	}

	return cacheImagesToRedis(ctx, variantID, images)
}

func ClearVariantImageCache(ctx context.Context, variantID string) error {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return err
	}

	if err := rdb.HDel(ctx, "product_variant_images", variantID).Err(); err != nil {
		return err
	}

	log.Printf("[CACHE] ✅ Cache variant %s dihapus dari Redis", variantID)
	return nil
}

func cacheImagesToRedis(ctx context.Context, variantID string, images []models.ProductVariantImage) error {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return err
	}

	key := "product_variant_images:" + variantID
	ttl := 300

	if len(images) == 0 {
		if err := rdb.Set(ctx, key, "0", time.Duration(ttl)*time.Second).Err(); err != nil {
			return err
		}
		log.Printf("[CACHE] ⚠️ Gambar kosong, simpan marker '0' dengan TTL untuk variant %s", variantID)
		return nil
	}

	var liteList []VariantImageLite
	for _, img := range images {
		liteList = append(liteList, VariantImageLite{
			ImageURL:       img.ImageURL,
			CloudImageID:   img.CloudImageID,
			IsVariantCover: img.IsVariantCover,
		})
	}

	data, err := json.Marshal(liteList)
	if err != nil {
		return err
	}

	if err := rdb.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err(); err != nil {
		return err
	}

	log.Printf("[CACHE] ✅ %d gambar variant %s disimpan di Redis (TTL %ds)", len(images), variantID, ttl)
	return nil
}

func getVariantImagesFromRedis(ctx context.Context, variantID string) ([]VariantImageLite, error) {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	key := "product_variant_images:" + variantID
	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if data == "0" {
		return []VariantImageLite{}, nil
	}

	var images []VariantImageLite
	if err := json.Unmarshal([]byte(data), &images); err != nil {
		return nil, err
	}

	return images, nil
}
