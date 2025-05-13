package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetVisitorProductDetail(c *fiber.Ctx) error {
	ctx := context.Background()

	summaries, err := fetcher.GetAllProductSummaries(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil daftar produk")
	}

	response := make([]fiber.Map, 0, len(summaries))
	for _, s := range summaries {
		response = append(response, fiber.Map{
			"id":                s.ID,
			"name":              s.Name,
			"slug":              s.Slug,
			"image_cover_url":   s.ImageCoverURL,
			"product_type_name": s.ProductTypeName,
			"product_form_name": s.ProductFormName,
			"min_price":         s.MinPrice,
			"default_variant": fiber.Map{
				"slug": s.DefaultVariant,
			},
		})
	}

	return c.JSON(response)
}
