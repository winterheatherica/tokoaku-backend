package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetHighlightedProductCards(c *fiber.Ctx) error {
	ctx := context.Background()

	var highlight models.HighlightedProduct
	if err := database.DB.Order("created_at DESC").First(&highlight).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Tidak ada produk highlight ditemukan")
	}

	product, err := fetcher.GetProductByID(ctx, highlight.ProductID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal ambil detail produk: "+err.Error())
	}

	variants, err := fetcher.GetProductVariantsByProductID(ctx, product.ID)
	if err != nil || len(variants) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Tidak ada varian untuk produk highlight ini")
	}

	var minPriceOriginal, maxPriceOriginal, minPriceDiscount, maxPriceDiscount *uint
	uniqueDiscounts := make(map[uint]fetcher.DiscountDTO)
	var allImages []string

	for _, variant := range variants {
		priceInfo, err := fetcher.GetPriceWithDiscountForUI(ctx, variant.ID)
		if err != nil || priceInfo.OriginalPrice == nil || priceInfo.FinalPrice == nil {
			continue
		}

		ori := *priceInfo.OriginalPrice
		dis := *priceInfo.FinalPrice

		if minPriceOriginal == nil || ori < *minPriceOriginal {
			tmp := ori
			minPriceOriginal = &tmp
		}
		if maxPriceOriginal == nil || ori > *maxPriceOriginal {
			tmp := ori
			maxPriceOriginal = &tmp
		}

		if minPriceDiscount == nil || dis < *minPriceDiscount {
			tmp := dis
			minPriceDiscount = &tmp
		}
		if maxPriceDiscount == nil || dis > *maxPriceDiscount {
			tmp := dis
			maxPriceDiscount = &tmp
		}

		for _, disc := range priceInfo.Discounts {
			uniqueDiscounts[disc.ID] = disc
		}

		variantImages, err := fetcher.GetAllVariantImages(ctx, variant.ID)
		if err == nil && len(variantImages) > 0 {
			for _, img := range variantImages {
				allImages = append(allImages, img.ImageURL)
			}
		}
	}

	allImages = append([]string{product.ImageCoverURL}, allImages...)

	var allDiscounts []fetcher.DiscountDTO
	for _, disc := range uniqueDiscounts {
		allDiscounts = append(allDiscounts, disc)
	}

	return c.JSON(fiber.Map{
		"product": fiber.Map{
			"id":              product.ID,
			"name":            product.Name,
			"slug":            product.Slug,
			"image_cover_url": product.ImageCoverURL,
		},
		"prices": fiber.Map{
			"min_original": minPriceOriginal,
			"max_original": maxPriceOriginal,
			"min_discount": minPriceDiscount,
			"max_discount": maxPriceDiscount,
		},
		"discounts": allDiscounts,
		"images":    allImages,
		"variants":  variants,
	})
}
