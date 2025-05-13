package seller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetProductVariants(c *fiber.Ctx) error {
	productID := c.Params("id")

	ctx := context.Background()
	variants, err := fetcher.GetProductVariantsByProductID(ctx, productID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch product variants")
	}

	response := make([]fiber.Map, 0, len(variants))
	for _, v := range variants {
		response = append(response, fiber.Map{
			"id":           v.ID,
			"variant_name": v.VariantName,
			"stock":        v.Stock,
		})
	}

	return c.JSON(response)
}
