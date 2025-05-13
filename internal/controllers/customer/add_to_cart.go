package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type AddToCartRequest struct {
	ProductSlug string `json:"product_slug"`
	VariantSlug string `json:"variant_slug"`
	Quantity    uint   `json:"quantity"`
}

func AddToCart(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)

	var req AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.Quantity == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Quantity must be at least 1")
	}

	var product models.Product
	if err := database.DB.Where("slug = ?", req.ProductSlug).First(&product).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	var variant models.ProductVariant
	if err := database.DB.Where("product_id = ? AND slug = ?", product.ID, req.VariantSlug).First(&variant).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Variant not found")
	}

	cartItem := models.Cart{
		CustomerID:       userUID,
		ProductVariantID: variant.ID,
		Quantity:         req.Quantity,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to add to cart")
	}

	return c.JSON(fiber.Map{
		"message": "Item added to cart",
		"cart":    cartItem,
	})
}
