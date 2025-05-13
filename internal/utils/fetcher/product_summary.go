package fetcher

import (
	"context"
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type ProductSummary struct {
	ID              string
	Name            string
	Slug            string
	ImageCoverURL   string
	ProductTypeName string
	ProductFormName string
	DefaultVariant  string
	MinPrice        *uint
}

func GetAllProductSummaries(ctx context.Context) ([]ProductSummary, error) {
	var products []models.Product

	err := database.DB.
		WithContext(ctx).
		Preload("ProductType").
		Preload("ProductForm").
		Preload("Variants.ProductVariantPrices").
		Order("created_at DESC").
		Find(&products).Error
	if err != nil {
		log.Printf("[DB] ❌ Gagal ambil produk: %v", err)
		return nil, err
	}

	summaries := make([]ProductSummary, 0, len(products))

	for _, product := range products {
		var defaultVariantSlug string
		var minPrice *uint = nil

		if len(product.Variants) > 0 {
			earliest := product.Variants[0]
			for _, v := range product.Variants {
				if v.CreatedAt.Before(earliest.CreatedAt) {
					earliest = v
				}
			}
			defaultVariantSlug = earliest.Slug

			for _, v := range product.Variants {
				for _, p := range v.ProductVariantPrices {
					if minPrice == nil || p.Price < *minPrice {
						price := p.Price
						minPrice = &price
					}
				}
			}
		}

		summaries = append(summaries, ProductSummary{
			ID:              product.ID,
			Name:            product.Name,
			Slug:            product.Slug,
			ImageCoverURL:   product.ImageCoverURL,
			ProductTypeName: product.ProductType.Name,
			ProductFormName: product.ProductForm.Form,
			DefaultVariant:  defaultVariantSlug,
			MinPrice:        minPrice,
		})
	}

	return summaries, nil
}
