package fetcher

import (
	"context"
	"fmt"
	"log"
	"strings"
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

type VariantPriceWithDiscount struct {
	ProductVariantID string            `json:"product_variant_id"`
	OriginalPrice    *uint             `json:"original_price"`
	FinalPrice       *uint             `json:"final_price"`
	AppliedDiscounts []models.Discount `json:"applied_discounts"`
	CreatedAt        *time.Time        `json:"created_at"`
}

type DiscountDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Value     uint   `json:"value"`
	ValueType string `json:"value_type"`
	Sponsor   string `json:"sponsor,omitempty"`
}

type PriceWithDiscountResponse struct {
	ProductVariantID string        `json:"product_variant_id"`
	OriginalPrice    *uint         `json:"original_price"`
	FinalPrice       *uint         `json:"final_price"`
	CreatedAt        *time.Time    `json:"created_at"`
	Discounts        []DiscountDTO `json:"discounts"`
}

func GetPriceWithDiscountForUI(ctx context.Context, variantID string) (*PriceWithDiscountResponse, error) {
	data, err := GetPriceWithDiscount(ctx, variantID)
	if err != nil {
		return nil, err
	}

	limit := 1
	if ce, err := GetCurrentEvent(ctx); err == nil && ce != nil {
		limit = int(ce.EventType.DiscountLimit)
	}

	var discounts []DiscountDTO
	for i, d := range data.AppliedDiscounts {
		if i >= limit {
			break
		}
		sponsorName := ""
		if d.DiscountSponsor.Role.ID != 0 {
			sponsorName = d.DiscountSponsor.Role.Name
		}

		discounts = append(discounts, DiscountDTO{
			ID:        d.ID,
			Name:      d.Name,
			Value:     d.Value,
			ValueType: d.ValueType.Name,
			Sponsor:   sponsorName,
		})
	}

	return &PriceWithDiscountResponse{
		ProductVariantID: data.ProductVariantID,
		OriginalPrice:    data.OriginalPrice,
		FinalPrice:       data.FinalPrice,
		CreatedAt:        data.CreatedAt,
		Discounts:        discounts,
	}, nil
}

func GetPriceWithDiscount(ctx context.Context, variantID string) (*VariantPriceWithDiscount, error) {
	priceData, err := GetLatestPriceForVariant(ctx, variantID)
	if err != nil || priceData.Price == nil {
		return &VariantPriceWithDiscount{
			ProductVariantID: variantID,
			OriginalPrice:    nil,
			FinalPrice:       nil,
			AppliedDiscounts: []models.Discount{},
			CreatedAt:        nil,
		}, nil
	}

	discounts, err := GetTopDiscountsByCurrentEventLimit(ctx, variantID)
	if err != nil {
		log.Printf("[DISCOUNT] ⚠️ Gagal ambil diskon aktif: %v", err)
		discounts = []models.Discount{}
	}

	price := float64(*priceData.Price)
	totalFlat := float64(0)
	totalPercentage := float64(1.0)
	var applied []models.Discount

	for _, d := range discounts {
		if d.Value == 0 {
			continue
		}
		switch strings.ToLower(d.ValueType.Name) {
		case "percentage":
			percentage := float64(d.Value) / 100.0
			totalPercentage *= (1 - percentage)
			applied = append(applied, d)
		}
	}

	for _, d := range discounts {
		if d.Value == 0 {
			continue
		}
		if strings.ToLower(d.ValueType.Name) == "flat" {
			totalFlat += float64(d.Value)
			applied = append(applied, d)
		}
	}

	afterPercent := price * totalPercentage
	afterFlat := afterPercent - totalFlat
	if afterFlat < 0 {
		afterFlat = 0
	}

	final := uint(afterFlat)

	log.Printf("[DEBUG] Harga asli: %.2f, setelah persen: %.2f, setelah flat: %.2f", price, afterPercent, afterFlat)

	priceUint := uint(price)
	return &VariantPriceWithDiscount{
		ProductVariantID: variantID,
		OriginalPrice:    &priceUint,
		FinalPrice:       &final,
		AppliedDiscounts: applied,
		CreatedAt:        priceData.CreatedAt,
	}, nil
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

	rdb, err := volatile.GetVolatileRedisClient(ctx)
	hashKey := fmt.Sprintf("product_variant:%s:%s", variant.ProductID, variantID)

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

	if err != nil {
		log.Printf("[DB] ❌ Harga variant %s tidak ditemukan: %v", variantID, err)
		return &VariantPriceResult{
			ProductVariantID: variantID,
			Price:            nil,
			CreatedAt:        nil,
		}, nil
	}

	result := &VariantPriceResult{
		ProductVariantID: price.ProductVariantID,
		Price:            &price.Price,
		CreatedAt:        &price.CreatedAt,
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

func GetPriceAtTimeForUI(ctx context.Context, variantID string, orderTime time.Time) (*PriceWithDiscountResponse, error) {
	var price models.ProductPrice
	err := database.DB.
		WithContext(ctx).
		Where("product_variant_id = ? AND created_at <= ?", variantID, orderTime).
		Order("created_at DESC").
		First(&price).Error

	if err != nil || price.Price == 0 {
		return &PriceWithDiscountResponse{
			ProductVariantID: variantID,
			OriginalPrice:    nil,
			FinalPrice:       nil,
			CreatedAt:        nil,
			Discounts:        []DiscountDTO{},
		}, nil
	}

	priceUint := price.Price
	createdAt := price.CreatedAt

	return &PriceWithDiscountResponse{
		ProductVariantID: variantID,
		OriginalPrice:    &priceUint,
		FinalPrice:       &priceUint,
		CreatedAt:        &createdAt,
		Discounts:        []DiscountDTO{},
	}, nil
}

func GetHistoricalPriceWithDiscount(ctx context.Context, variantID string, orderTime time.Time) (*PriceWithDiscountResponse, error) {
	var price models.ProductPrice
	err := database.DB.
		WithContext(ctx).
		Where("product_variant_id = ? AND created_at <= ?", variantID, orderTime).
		Order("created_at DESC").
		First(&price).Error

	if err != nil || price.Price == 0 {
		log.Printf("[HISTORIC_PRICE] Tidak ditemukan harga untuk %s sebelum %v", variantID, orderTime)
		return &PriceWithDiscountResponse{
			ProductVariantID: variantID,
			OriginalPrice:    nil,
			FinalPrice:       nil,
			CreatedAt:        nil,
			Discounts:        []DiscountDTO{},
		}, nil
	}

	discounts, err := GetTopDiscountsByOrderTimestamp(ctx, variantID, orderTime)
	if err != nil {
		log.Printf("[HISTORIC_DISCOUNT] Gagal ambil diskon historis: %v", err)
		discounts = []models.Discount{}
	}

	priceFloat := float64(price.Price)
	totalPercentage := 1.0
	totalFlat := float64(0)
	appliedDiscounts := []DiscountDTO{}

	for _, d := range discounts {
		if d.Value == 0 {
			continue
		}
		switch strings.ToLower(d.ValueType.Name) {
		case "percentage":
			pct := float64(d.Value) / 100.0
			totalPercentage *= (1 - pct)
		case "flat":
			totalFlat += float64(d.Value)
		}

		sponsor := ""
		if d.DiscountSponsor.Role.ID != 0 {
			sponsor = d.DiscountSponsor.Role.Name
		}

		appliedDiscounts = append(appliedDiscounts, DiscountDTO{
			ID:        d.ID,
			Name:      d.Name,
			Value:     d.Value,
			ValueType: d.ValueType.Name,
			Sponsor:   sponsor,
		})
	}

	afterPercentage := priceFloat * totalPercentage
	afterFlat := afterPercentage - totalFlat
	if afterFlat < 0 {
		afterFlat = 0
	}
	finalPrice := uint(afterFlat)

	log.Printf("[HISTORIC_PRICE] Original: %d, Final after discounts: %d (%% multiplier=%.4f, flat=%.2f)",
		price.Price, finalPrice, totalPercentage, totalFlat)

	originalPriceUint := uint(price.Price)
	return &PriceWithDiscountResponse{
		ProductVariantID: variantID,
		OriginalPrice:    &originalPriceUint,
		FinalPrice:       &finalPrice,
		CreatedAt:        &price.CreatedAt,
		Discounts:        appliedDiscounts,
	}, nil
}
