package fetcher

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetTopDiscountsByCurrentEventLimit(ctx context.Context, variantID string) ([]models.Discount, error) {
	var variant models.ProductVariant
	if err := database.DB.WithContext(ctx).
		Preload("Product.ProductCategories.Category").
		Preload("Product.ProductType").
		Where("id = ?", variantID).
		First(&variant).Error; err != nil {
		return nil, err
	}

	product := variant.Product
	productTypeID := product.ProductType.ID

	log.Printf("[DEBUG] Variant %s milik Product %s (TypeID: %d), total kategori: %d",
		variant.ID, product.ID, productTypeID, len(product.ProductCategories))

	var categoryIDs []uint
	for _, pc := range product.ProductCategories {
		categoryIDs = append(categoryIDs, pc.CategoryID)
	}

	var discountIDs []uint

	var variantDiscounts []models.ProductVariantDiscount
	database.DB.WithContext(ctx).
		Where("product_variant_id = ?", variant.ID).
		Find(&variantDiscounts)
	log.Printf("[DEBUG] Jumlah diskon dari variant pivot: %d", len(variantDiscounts))
	for _, pd := range variantDiscounts {
		discountIDs = append(discountIDs, pd.DiscountID)
	}

	var typeDiscounts []models.ProductTypeDiscount
	database.DB.WithContext(ctx).
		Where("product_type_id = ?", productTypeID).
		Find(&typeDiscounts)
	log.Printf("[DEBUG] Jumlah diskon dari type pivot: %d", len(typeDiscounts))
	for _, td := range typeDiscounts {
		discountIDs = append(discountIDs, td.DiscountID)
	}

	var categoryDiscounts []models.CategoryDiscount
	if len(categoryIDs) > 0 {
		database.DB.WithContext(ctx).
			Where("category_id IN ?", categoryIDs).
			Find(&categoryDiscounts)
		log.Printf("[DEBUG] Jumlah diskon dari category pivot: %d", len(categoryDiscounts))
		for _, cd := range categoryDiscounts {
			discountIDs = append(discountIDs, cd.DiscountID)
		}
	}

	discountIDMap := make(map[uint]struct{})
	var uniqueIDs []uint
	for _, id := range discountIDs {
		if _, exists := discountIDMap[id]; !exists {
			discountIDMap[id] = struct{}{}
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	log.Printf("[DEBUG] Total diskon unik setelah merge: %d", len(uniqueIDs))

	if len(uniqueIDs) == 0 {
		log.Println("[FETCH] ❌ Tidak ada ID diskon ditemukan dari semua pivot")
		return []models.Discount{}, nil
	}

	var discounts []models.Discount
	if err := database.DB.WithContext(ctx).
		Preload("ValueType").
		Preload("DiscountSponsor.Role").
		Where("id IN ?", uniqueIDs).
		Where(`"start_at" <= ? AND "end_at" >= ?`, time.Now(), time.Now()).
		Order("value_type_id ASC, value DESC").
		Limit(3).
		Find(&discounts).Error; err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Final jumlah diskon aktif yang diambil: %d", len(discounts))
	totalFlat := 0
	totalPercent := 0
	for _, d := range discounts {
		switch strings.ToLower(d.ValueType.Name) {
		case "flat":
			totalFlat++
		case "percentage":
			totalPercent++
		}
		log.Printf("  ↳ Diskon: %s (ID: %d) - Value: %d %s", d.Name, d.ID, d.Value, d.ValueType.Name)
	}
	log.Printf("[FETCH] ✅ Ambil %d diskon (Flat %d + Percentage %d) dari pivot", len(discounts), totalFlat, totalPercent)

	return discounts, nil
}

func GetTopDiscountsByOrderTimestamp(ctx context.Context, variantID string, orderTime time.Time) ([]models.Discount, error) {
	var variant models.ProductVariant
	if err := database.DB.WithContext(ctx).
		Preload("Product.ProductCategories.Category").
		Preload("Product.ProductType").
		Where("id = ?", variantID).
		First(&variant).Error; err != nil {
		return nil, err
	}

	product := variant.Product
	productTypeID := product.ProductType.ID

	var categoryIDs []uint
	for _, pc := range product.ProductCategories {
		categoryIDs = append(categoryIDs, pc.CategoryID)
	}

	var discountIDs []uint

	var variantDiscounts []models.ProductVariantDiscount
	database.DB.WithContext(ctx).
		Where("product_variant_id = ?", variant.ID).
		Find(&variantDiscounts)
	for _, pd := range variantDiscounts {
		discountIDs = append(discountIDs, pd.DiscountID)
	}

	var typeDiscounts []models.ProductTypeDiscount
	database.DB.WithContext(ctx).
		Where("product_type_id = ?", productTypeID).
		Find(&typeDiscounts)
	for _, td := range typeDiscounts {
		discountIDs = append(discountIDs, td.DiscountID)
	}

	var categoryDiscounts []models.CategoryDiscount
	if len(categoryIDs) > 0 {
		database.DB.WithContext(ctx).
			Where("category_id IN ?", categoryIDs).
			Find(&categoryDiscounts)
		for _, cd := range categoryDiscounts {
			discountIDs = append(discountIDs, cd.DiscountID)
		}
	}

	discountIDMap := make(map[uint]struct{})
	var uniqueIDs []uint
	for _, id := range discountIDs {
		if _, exists := discountIDMap[id]; !exists {
			discountIDMap[id] = struct{}{}
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	if len(uniqueIDs) == 0 {
		log.Println("[FETCH] ❌ Tidak ada ID diskon ditemukan dari semua pivot")
		return []models.Discount{}, nil
	}

	var event models.CurrentEvent
	if err := database.DB.WithContext(ctx).
		Preload("EventType").
		Where("start_at <= ? AND end_at >= ?", orderTime, orderTime).
		Order("start_at DESC").
		First(&event).Error; err != nil {
		log.Printf("[EVENT] ⚠️ Event tidak ditemukan pada %v: %v", orderTime, err)
		return []models.Discount{}, nil
	}

	limit := int(event.EventType.DiscountLimit)
	if limit <= 1 {
		limit = 3
	}

	var discounts []models.Discount
	if err := database.DB.WithContext(ctx).
		Preload("ValueType").
		Preload("DiscountSponsor.Role").
		Where("id IN ?", uniqueIDs).
		Where("start_at <= ? AND end_at >= ?", orderTime, orderTime).
		Order("value_type_id ASC, value DESC").
		Limit(limit).
		Find(&discounts).Error; err != nil {
		return nil, err
	}

	totalFlat := 0
	totalPercent := 0
	for _, d := range discounts {
		switch strings.ToLower(d.ValueType.Name) {
		case "flat":
			totalFlat++
		case "percentage":
			totalPercent++
		}
		log.Printf("  ↳ Diskon: %s (ID: %d) - Value: %d %s", d.Name, d.ID, d.Value, d.ValueType.Name)
	}
	log.Printf("[FETCH-HISTORIC] ✅ Ambil %d diskon (Flat %d + Percentage %d) dari event pada %v", len(discounts), totalFlat, totalPercent, orderTime)

	return discounts, nil
}

func GetEventByTimestamp(ctx context.Context, t time.Time) (*models.CurrentEvent, error) {
	var event models.CurrentEvent
	err := database.DB.WithContext(ctx).
		Preload("EventType").
		Where("start_at <= ? AND end_at >= ?", t, t).
		Order("start_at DESC").
		First(&event).Error

	if err != nil {
		return nil, err
	}
	return &event, nil
}
