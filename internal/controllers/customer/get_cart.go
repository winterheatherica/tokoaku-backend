package customer

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetGroupedCart(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)

	var cartItems []models.Cart
	if err := database.DB.
		Preload("ProductVariant.Product.Seller").
		Preload("ProductVariant.Product").
		Preload("ProductVariant").
		Where("customer_id = ?", userUID).
		Order("created_at DESC").
		Find(&cartItems).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil keranjang")
	}

	sellerGrouped := make(map[string]map[string][]fiber.Map)

	for _, item := range cartItems {
		variantImages, err := fetcher.GetAllVariantImages(context.Background(), item.ProductVariantID)
		if err != nil {
			variantImages = []models.ProductVariantImage{}
		}

		imageURL := item.ProductVariant.Product.ImageCoverURL
		for _, img := range variantImages {
			if img.IsVariantCover {
				imageURL = img.ImageURL
				break
			}
		}

		images := []string{}
		for _, img := range variantImages {
			images = append(images, img.ImageURL)
		}

		sellerName := "-"
		if item.ProductVariant.Product.Seller.Name != nil {
			sellerName = *item.ProductVariant.Product.Seller.Name
		}
		productSlug := item.ProductVariant.Product.Slug

		cartData := fiber.Map{
			"product_variant_id": item.ProductVariantID,
			"product_name":       item.ProductVariant.Product.Name,
			"product_slug":       productSlug,
			"variant_name":       item.ProductVariant.VariantName,
			"variant_slug":       item.ProductVariant.Slug,
			"quantity":           item.Quantity,
			"image_url":          imageURL,
			"variant_images":     images,
			"added_at":           item.CreatedAt,
			"is_selected":        item.IsSelected,
		}

		if _, ok := sellerGrouped[sellerName]; !ok {
			sellerGrouped[sellerName] = make(map[string][]fiber.Map)
		}
		sellerGrouped[sellerName][productSlug] = append(sellerGrouped[sellerName][productSlug], cartData)
	}

	var response []fiber.Map
	for sellerName, products := range sellerGrouped {
		var groupedProducts []fiber.Map
		for _, items := range products {
			productName := ""
			if len(items) > 0 {
				if name, ok := items[0]["product_name"].(string); ok {
					productName = name
				}
			}

			groupedProducts = append(groupedProducts, fiber.Map{
				"product_name": productName,
				"cart_items":   items,
			})
		}

		response = append(response, fiber.Map{
			"seller_name": sellerName,
			"products":    groupedProducts,
		})
	}

	return c.JSON(response)
}
