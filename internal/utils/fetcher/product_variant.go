package fetcher

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func GetProductVariantsByProductID(ctx context.Context, productID string) ([]models.ProductVariant, error) {
	if variants, ok := getVariantsFromCache(ctx, productID); ok {
		return variants, nil
	}

	var variants []models.ProductVariant
	if err := database.DB.
		Where("product_id = ?", productID).
		Find(&variants).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil variant produk %s: %v", productID, err)
		return nil, err
	}

	log.Printf("[DB] ✅ Ambil %d variant produk %s dari DB", len(variants), productID)
	cacheVariants(ctx, productID, variants)
	return variants, nil
}

func GetProductVariantBySlug(ctx context.Context, productID, variantSlug string) (*models.ProductVariant, error) {
	if variant, ok := getVariantBySlugFromCache(ctx, productID, variantSlug); ok {
		return variant, nil
	}

	var variant models.ProductVariant
	if err := database.DB.
		WithContext(ctx).
		Where("product_id = ? AND slug = ?", productID, variantSlug).
		First(&variant).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil variant %s untuk produk %s: %v", variantSlug, productID, err)
		return nil, err
	}

	log.Printf("[DB] ✅ Variant %s untuk produk %s berhasil diambil dari DB", variantSlug, productID)
	CacheSingleVariant(ctx, productID, &variant)
	return &variant, nil
}

func getVariantsFromCache(ctx context.Context, productID string) ([]models.ProductVariant, bool) {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return nil, false
	}

	setKey := fmt.Sprintf("product_variant:%s:ids", productID)
	ids, err := rdb.SMembers(ctx, setKey).Result()
	if err != nil || len(ids) == 0 {
		return nil, false
	}

	var variants []models.ProductVariant
	for _, id := range ids {
		hashKey := fmt.Sprintf("product_variant:%s:%s", productID, id)
		data, err := rdb.HGetAll(ctx, hashKey).Result()
		if err != nil || len(data) == 0 {
			continue
		}

		if variant := parseVariantHash(data); variant != nil {
			variants = append(variants, *variant)
		}
	}

	log.Printf("[CACHE] ✅ Ambil %d variant produk %s dari Redis", len(variants), productID)
	return variants, true
}

func getVariantBySlugFromCache(ctx context.Context, productID, slug string) (*models.ProductVariant, bool) {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return nil, false
	}

	setKey := fmt.Sprintf("product_variant:%s:ids", productID)
	ids, err := rdb.SMembers(ctx, setKey).Result()
	if err != nil || len(ids) == 0 {
		return nil, false
	}

	for _, id := range ids {
		hashKey := fmt.Sprintf("product_variant:%s:%s", productID, id)
		data, err := rdb.HGetAll(ctx, hashKey).Result()
		if err != nil || len(data) == 0 {
			continue
		}

		if data["slug"] == slug {
			if variant := parseVariantHash(data); variant != nil {
				log.Printf("[CACHE] ✅ Variant %s untuk produk %s ditemukan di Redis", slug, productID)
				return variant, true
			}
		}
	}

	return nil, false
}

func cacheVariants(ctx context.Context, productID string, variants []models.ProductVariant) {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return
	}

	setKey := fmt.Sprintf("product_variant:%s:ids", productID)
	var ids []string

	for _, v := range variants {
		hashKey := fmt.Sprintf("product_variant:%s:%s", productID, v.ID)
		data := map[string]interface{}{
			"id":           v.ID,
			"variant_name": v.VariantName,
			"slug":         v.Slug,
			"stock":        v.Stock,
			"created_at":   v.CreatedAt.Format(time.RFC3339),
		}
		if err := rdb.HSet(ctx, hashKey, data).Err(); err == nil {
			_ = rdb.Expire(ctx, hashKey, 5*time.Minute).Err()
			ids = append(ids, v.ID)
		} else {
			log.Printf("[CACHE] ❌ Gagal simpan variant %s ke Redis: %v", v.ID, err)
		}
	}

	if len(ids) > 0 {
		_ = rdb.SAdd(ctx, setKey, ids).Err()
		_ = rdb.Expire(ctx, setKey, 5*time.Minute).Err()
		log.Printf("[CACHE] ✅ Simpan %d variant produk %s ke Redis", len(ids), productID)
	}
}

func CacheSingleVariant(ctx context.Context, productID string, variant *models.ProductVariant) error {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return err
	}

	hashKey := fmt.Sprintf("product_variant:%s:%s", productID, variant.ID)
	data := map[string]interface{}{
		"id":           variant.ID,
		"variant_name": variant.VariantName,
		"slug":         variant.Slug,
		"stock":        variant.Stock,
		"created_at":   variant.CreatedAt.Format(time.RFC3339),
	}
	if err := rdb.HSet(ctx, hashKey, data).Err(); err != nil {
		return err
	}

	if err := rdb.Expire(ctx, hashKey, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func parseVariantHash(data map[string]string) *models.ProductVariant {
	createdAt, err := time.Parse(time.RFC3339, data["created_at"])
	if err != nil {
		return nil
	}
	return &models.ProductVariant{
		ID:          data["id"],
		VariantName: data["variant_name"],
		Slug:        data["slug"],
		Stock:       parseUint(data["stock"]),
		CreatedAt:   createdAt,
	}
}

func parseUint(s string) uint {
	n, _ := strconv.ParseUint(s, 10, 32)
	return uint(n)
}

func InvalidateVariantCache(ctx context.Context, productID string) error {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product_variant:%s:ids", productID)
	return rdb.Del(ctx, key).Err()
}
