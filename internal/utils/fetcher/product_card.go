package fetcher

import (
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type ProductCard struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	ImageCoverURL   string `json:"image_cover_url"`
	ProductTypeName string `json:"product_type_name"`
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

		var prices []models.ProductPrice
		if len(variants) > 0 {
			for _, v := range variants {
				var latestPrice models.ProductPrice
				if err := database.DB.
					Where("product_variant_id = ?", v.ID).
					Order("created_at DESC").
					First(&latestPrice).Error; err == nil {
					prices = append(prices, latestPrice)
				}
			}
		}

		var minPrice *uint
		var maxPrice *uint
		for _, p := range prices {
			if minPrice == nil || p.Price < *minPrice {
				v := p.Price
				minPrice = &v
			}
			if maxPrice == nil || p.Price > *maxPrice {
				v := p.Price
				maxPrice = &v
			}
		}

		response = append(response, ProductCard{
			ID:              product.ID,
			Name:            product.Name,
			Slug:            product.Slug,
			ImageCoverURL:   product.ImageCoverURL,
			ProductTypeName: product.ProductType.Name,
			ProductFormName: product.ProductForm.Form,
			VariantSlug:     variantSlug,
			MinPrice:        minPrice,
			MaxPrice:        maxPrice,
		})
	}

	return response, nil
}
