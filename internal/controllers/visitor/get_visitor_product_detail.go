package visitor

import (
	"context"
	"log"

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
		var latestPrice uint = 0

		if s.DefaultVariant != "" {
			priceWithDiscount, err := fetcher.GetPriceWithDiscountForUI(ctx, s.DefaultVariant)
			if err != nil {
				log.Printf("[ERROR] Gagal ambil harga diskon untuk variant %s: %v", s.DefaultVariant, err)
			} else if priceWithDiscount != nil && priceWithDiscount.FinalPrice != nil {
				latestPrice = *priceWithDiscount.FinalPrice

				if len(priceWithDiscount.Discounts) > 0 {
					for _, d := range priceWithDiscount.Discounts {
						log.Printf("[DISCOUNT] Produk %s pakai diskon ID %d - %s (%s)", s.Name, d.ID, d.Name, d.ValueType)
					}
				} else {
					log.Printf("[DISCOUNT] Produk %s tidak pakai diskon aktif", s.Name)
				}
			}
		}

		response = append(response, fiber.Map{
			"id":                s.ID,
			"name":              s.Name,
			"slug":              s.Slug,
			"image_cover_url":   s.ImageCoverURL,
			"product_type_name": s.ProductTypeName,
			"product_form_name": s.ProductFormName,
			"latest_price":      latestPrice,
			"default_variant": fiber.Map{
				"slug": s.DefaultVariant,
			},
		})
	}

	return c.JSON(response)
}
