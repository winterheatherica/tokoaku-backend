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

type VariantPriceResult struct {
	ProductVariantID string     `json:"product_variant_id"`
	Price            *uint      `json:"price"`
	CreatedAt        *time.Time `json:"created_at"`
}

func GetLatestPriceForVariant(ctx context.Context, variantID string) (*VariantPriceResult, error) {
	var variant models.ProductVariant
	if err := database.DB.
		WithContext(ctx).
		Select("product_id").
		Where("id = ?", variantID).
		First(&variant).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil product_id dari variant %s: %v", variantID, err)
		return nil, err
	}
	productID := variant.ProductID

	rdb, err := volatile.GetVolatileRedisClient(ctx)
	hashKey := fmt.Sprintf("product_variant:%s:%s", productID, variantID)

	if err == nil {
		data, err := rdb.HMGet(ctx, hashKey, "latest_price", "latest_price_created_at").Result()
		if err == nil && len(data) == 2 && data[0] != nil && data[1] != nil {
			priceUint, priceErr := parseUintFromString(data[0])
			createdAt, timeErr := time.Parse(time.RFC3339, fmt.Sprint(data[1]))

			if priceErr == nil && timeErr == nil {
				log.Printf("[CACHE] ✅ Harga variant %s ditemukan di Redis hash", variantID)
				return &VariantPriceResult{
					ProductVariantID: variantID,
					Price:            &priceUint,
					CreatedAt:        &createdAt,
				}, nil
			}
			log.Printf("[CACHE] ⚠️ Gagal decode harga variant %s: %v %v", variantID, priceErr, timeErr)
		}
	}

	var price models.ProductPrice
	err = database.DB.
		Where("product_variant_id = ?", variantID).
		Order("created_at DESC").
		First(&price).Error

	var result *VariantPriceResult

	if err != nil {
		log.Printf("[DB] ❌ Harga variant %s tidak ditemukan: %v", variantID, err)
		result = &VariantPriceResult{
			ProductVariantID: variantID,
			Price:            nil,
			CreatedAt:        nil,
		}
	} else {
		result = &VariantPriceResult{
			ProductVariantID: price.ProductVariantID,
			Price:            &price.Price,
			CreatedAt:        &price.CreatedAt,
		}
	}

	if rdb != nil && result.Price != nil && result.CreatedAt != nil {
		_ = rdb.HSet(ctx, hashKey, map[string]interface{}{
			"latest_price":            *result.Price,
			"latest_price_created_at": result.CreatedAt.Format(time.RFC3339),
		}).Err()
		_ = rdb.Expire(ctx, hashKey, 5*time.Minute).Err()
		log.Printf("[CACHE] ✅ Harga variant %s disimpan ke Redis hash utama", variantID)
	}

	return result, nil
}

func parseUintFromString(val interface{}) (uint, error) {
	str := fmt.Sprintf("%v", val)
	n, err := strconv.ParseUint(str, 10, 32)
	return uint(n), err
}
