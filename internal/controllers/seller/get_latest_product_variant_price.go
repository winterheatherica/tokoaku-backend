package seller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetLatestProductVariantPrice(c *fiber.Ctx) error {
	variantID := c.Params("id")
	if variantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing variant ID")
	}

	ctx := context.Background()

	var variant models.ProductVariant
	if err := database.DB.
		WithContext(ctx).
		Select("product_id").
		Where("id = ?", variantID).
		First(&variant).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Variant not found")
	}

	result, err := fetcher.GetLatestPriceForVariant(ctx, variantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve price")
	}

	return c.JSON(fiber.Map{
		"product_variant_id": result.ProductVariantID,
		"price":              result.Price,
		"created_at":         result.CreatedAt,
	})
}
