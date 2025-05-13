package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetCart(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)

	var cartItems []models.Cart
	if err := database.DB.
		Preload("ProductVariant.Product").
		Preload("ProductVariant").
		Where("customer_id = ?", userUID).
		Order("created_at DESC").
		Find(&cartItems).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil keranjang")
	}

	var response []fiber.Map
	for _, item := range cartItems {
		var coverImage models.ProductVariantImage
		database.DB.
			Where("product_variant_id = ? AND is_variant_cover = true", item.ProductVariantID).
			First(&coverImage)

		response = append(response, fiber.Map{
			"product_name": item.ProductVariant.Product.Name,
			"product_slug": item.ProductVariant.Product.Slug,
			"variant_name": item.ProductVariant.VariantName,
			"variant_slug": item.ProductVariant.Slug,
			"quantity":     item.Quantity,
			"image_url":    coverImage.ImageURL,
			"added_at":     item.CreatedAt,
		})
	}

	return c.JSON(response)
}
