package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetProductByType(c *fiber.Ctx) error {
	typeID := c.Query("type_id")
	var products []models.Product

	query := database.DB.Preload("ProductType").Preload("ProductType")

	if typeID != "" {
		query = query.Where("product_type_id = ?", typeID)
	}

	if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var response []fetcher.ProductCard
	ctx := context.Background()

	for _, product := range products {
		var variants []models.ProductVariant
		if err := database.DB.Where("product_id = ?", product.ID).Order("created_at ASC").Find(&variants).Error; err != nil {
			continue
		}

		var variantSlug string
		if len(variants) > 0 {
			variantSlug = variants[0].Slug
		}

		var minPrice *uint
		var maxPrice *uint

		for _, v := range variants {
			priceWithDiscount, err := fetcher.GetPriceWithDiscountForUI(ctx, v.ID)
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
		}

		response = append(response, fetcher.ProductCard{
			ID:              product.ID,
			Name:            product.Name,
			Slug:            product.Slug,
			ImageCoverURL:   product.ImageCoverURL,
			ProductTypeID:   product.ProductTypeID,
			ProductTypeName: product.ProductType.Name,
			VariantSlug:     variantSlug,
			MinPrice:        minPrice,
			MaxPrice:        maxPrice,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil daftar produk",
		"data":    response,
	})
}
