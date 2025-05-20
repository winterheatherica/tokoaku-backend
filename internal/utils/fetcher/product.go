package fetcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal koneksi ke Redis: %v", err)
		return nil, err
	}

	productID, err := rdb.HGet(ctx, "product_slug_map", slug).Result()
	if err != nil || productID == "" {
		log.Printf("[CACHE] ℹ️ Slug %s tidak ditemukan di Redis, ambil dari DB", slug)
		return fetchProductFromDBAndCache(ctx, slug)
	}

	log.Printf("[CACHE] 🔁 Slug %s → ProductID %s ditemukan di Redis", slug, productID)

	metaKey := fmt.Sprintf("product:%s", productID)
	metaData, err := rdb.HGetAll(ctx, metaKey).Result()
	if err != nil || len(metaData) == 0 {
		log.Printf("[CACHE] ⚠️ Metadata produk %s kosong di Redis, fallback ke DB", productID)
		return fetchProductFromDBAndCache(ctx, slug)
	}

	product := &models.Product{
		ID:            metaData["id"],
		Name:          metaData["name"],
		Description:   metaData["description"],
		Slug:          metaData["slug"],
		ImageCoverURL: metaData["image_cover_url"],
		ProductType:   models.ProductType{Name: metaData["product_type"]},
		ProductForm:   models.ProductForm{Form: metaData["product_form"]},
	}

	variants, _ := getVariantsFromCache(ctx, productID)
	product.Variants = variants

	log.Printf("[CACHE] ✅ Product %s diambil dari Redis secara modular", slug)
	return product, nil
}

func fetchProductFromDBAndCache(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	if err := database.DB.
		WithContext(ctx).
		Preload("ProductType").
		Preload("ProductForm").
		Preload("CloudService").
		Preload("Variants").
		Where("slug = ?", slug).
		First(&product).Error; err != nil {
		log.Printf("[DB] ❌ Gagal ambil product %s: %v", slug, err)
		return nil, err
	}

	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err == nil {
		if err := rdb.HSet(ctx, "product_slug_map", slug, product.ID).Err(); err != nil {
			log.Printf("[CACHE] ⚠️ Gagal simpan slug map: %v", err)
		}

		metaKey := fmt.Sprintf("product:%s", product.ID)
		data := map[string]interface{}{
			"id":              product.ID,
			"name":            product.Name,
			"description":     product.Description,
			"slug":            product.Slug,
			"seller_id":       product.SellerID,
			"product_type":    product.ProductType.Name,
			"product_form":    product.ProductForm.Form,
			"product_type_id": fmt.Sprintf("%d", product.ProductTypeID),
			"image_cover_url": product.ImageCoverURL,
			"cloud_image_id":  fmt.Sprintf("%d", product.CloudImageID),
			"product_form_id": fmt.Sprintf("%d", product.ProductFormID),
			"created_at":      product.CreatedAt.Format(time.RFC3339),
			"updated_at":      product.UpdatedAt.Format(time.RFC3339),
		}

		if err := rdb.HSet(ctx, metaKey, data).Err(); err != nil {
			log.Printf("[CACHE] ⚠️ Gagal simpan metadata product ke Redis: %v", err)
		}

		if err := rdb.Expire(ctx, metaKey, 1*time.Hour).Err(); err != nil {
			log.Printf("[CACHE] ⚠️ Gagal set TTL untuk product %s: %v", product.ID, err)
		}
	} else {
		log.Printf("[CACHE] ⚠️ Gagal koneksi ke Redis saat cache product %s: %v", product.ID, err)
	}

	cacheVariants(ctx, product.ID, product.Variants)

	return &product, nil
}
