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
		images = []models.ProductVariantImage{}
	}

	latestPrice, err := fetcher.GetLatestPriceForVariant(ctx, variant.ID)
	if err != nil || latestPrice.Price == nil {
		zero := uint(0)
		latestPrice.Price = &zero
	}

	allVariants, err := fetcher.GetProductVariantsByProductID(ctx, product.ID)
	if err != nil {
		allVariants = []models.ProductVariant{}
	}

	if len(images) == 0 && product.ImageCoverURL != "" {
		images = []models.ProductVariantImage{
			{
				ImageURL:       product.ImageCoverURL,
				IsVariantCover: true,
			},
		}
	}

	return c.JSON(fiber.Map{
		"product": fiber.Map{
			"id":              product.ID,
			"name":            product.Name,
			"description":     product.Description,
			"image_cover_url": product.ImageCoverURL,
			"product_type":    product.ProductType.Name,
			"product_form":    product.ProductForm.Form,
			"slug":            product.Slug,
		},
		"variant": fiber.Map{
			"id":           variant.ID,
			"name":         variant.VariantName,
			"slug":         variant.Slug,
			"stock":        variant.Stock,
			"images":       images,
			"latest_price": latestPrice.Price,
		},
		"all_variants": func() []fiber.Map {
			list := []fiber.Map{}
			for _, v := range allVariants {
				list = append(list, fiber.Map{
					"id":           v.ID,
					"variant_name": v.VariantName,
					"slug":         v.Slug,
					"created_at":   v.CreatedAt,
				})
			}
			return list
		}(),
	})
}
