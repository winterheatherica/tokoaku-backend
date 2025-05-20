package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetVisitorProductDetailByVariant(c *fiber.Ctx) error {
	productSlug := c.Params("productSlug")
	variantSlug := c.Params("variantSlug")
	ctx := context.Background()

	product, err := fetcher.GetProductBySlug(ctx, productSlug)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	variant, err := fetcher.GetProductVariantBySlug(ctx, product.ID, variantSlug)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	images, err := fetcher.GetAllVariantImages(ctx, variant.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	priceInfo, err := fetcher.GetPriceWithDiscountForUI(ctx, variant.ID)
	if err != nil || priceInfo.OriginalPrice == nil || priceInfo.FinalPrice == nil {
		zero := uint(0)
		priceInfo = &fetcher.PriceWithDiscountResponse{
			ProductVariantID: variant.ID,
			OriginalPrice:    &zero,
			FinalPrice:       &zero,
			Discounts:        []fetcher.DiscountDTO{},
		}
	}

	allVariants, err := fetcher.GetProductVariantsByProductID(ctx, product.ID)
	if err != nil {
		allVariants = []models.ProductVariant{}
	}

	return c.JSON(fiber.Map{
		"product": product,
		"variant": fiber.Map{
			"id":             variant.ID,
			"name":           variant.VariantName,
			"slug":           variant.Slug,
			"stock":          variant.Stock,
			"created_at":     variant.CreatedAt,
			"original_price": priceInfo.OriginalPrice,
			"final_price":    priceInfo.FinalPrice,
			"discounts":      priceInfo.Discounts,
			"images":         images,
		},
		"all_variants": allVariants,
	})
}
