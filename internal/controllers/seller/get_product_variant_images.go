package seller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetProductVariantImages(c *fiber.Ctx) error {
	variantID := c.Params("id")
	if variantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Variant ID is required")
	}

	images, err := fetcher.GetAllVariantImages(c.Context(), variantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch variant images")
	}

	response := make([]fiber.Map, 0, len(images))
	for _, img := range images {
		response = append(response, fiber.Map{
			"image_url":        img.ImageURL,
			"is_variant_cover": img.IsVariantCover,
		})
	}

	return c.JSON(response)
}
