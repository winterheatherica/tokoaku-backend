package fetcher

import (
	"context"
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type ProductCard struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	ImageCoverURL   string `json:"image_cover_url"`
	ProductTypeName string `json:"product_type_name"`
	ProductTypeID   uint   `json:"product_type_id"`
	ProductFormID   uint   `json:"product_form_id"`
	ProductFormName string `json:"product_form_name"`
	VariantSlug     string `json:"variant_slug"`
	MinPrice        *uint  `json:"min_price"`
	MaxPrice        *uint  `json:"max_price"`
}

func FetchProductCardList() ([]ProductCard, error) {
	var products []models.Product
	if err := database.DB.
		Preload("ProductType").
		Preload("ProductForm").
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return nil, err
	}

	ctx := context.Background()
	var response []ProductCard

	for _, product := range products {
		var variants []models.ProductVariant
		if err := database.DB.
			Where("product_id = ?", product.ID).
			Order("created_at ASC").
			Find(&variants).Error; err != nil {
			continue
		}

		var variantSlug string
		if len(variants) > 0 {
			variantSlug = variants[0].Slug
		}

		var minPrice *uint
		var maxPrice *uint

		for _, v := range variants {
			priceWithDiscount, err := GetPriceWithDiscountForUI(ctx, v.ID)
			if err != nil || priceWithDiscount == nil || priceWithDiscount.FinalPrice == nil {
				continue
			}

			final := *priceWithDiscount.FinalPrice

			if minPrice == nil || final < *minPrice {
				tmp := final
				minPrice = &tmp
			}
			if maxPrice == nil || final > *maxPrice {
				tmp := final
				maxPrice = &tmp
			}

			for _, d := range priceWithDiscount.Discounts {
				log.Printf("[DISCOUNT] Produk %s pakai diskon ID %d - %s (%s)", product.Name, d.ID, d.Name, d.ValueType)
			}
		}

		response = append(response, ProductCard{
			ID:              product.ID,
			Name:            product.Name,
			Slug:            product.Slug,
			ImageCoverURL:   product.ImageCoverURL,
			ProductTypeName: product.ProductType.Name,
			ProductFormID:   product.ProductFormID,
			ProductFormName: product.ProductForm.Form,
			VariantSlug:     variantSlug,
			MinPrice:        minPrice,
			MaxPrice:        maxPrice,
		})

	}

	return response, nil
}
