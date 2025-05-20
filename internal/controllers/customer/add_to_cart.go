package customer

import (
	"log"

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

	var existingCart models.Cart
	err := database.DB.
		Where("customer_id = ? AND product_variant_id = ? AND is_converted = false", userUID, variant.ID).
		First(&existingCart).Error

	if err == nil {
		newQty := existingCart.Quantity + req.Quantity
		if err := database.DB.
			Model(&models.Cart{}).
			Where("customer_id = ? AND product_variant_id = ? AND is_converted = false", userUID, variant.ID).
			Update("quantity", newQty).Error; err != nil {
			log.Println("[ERROR] Gagal update quantity:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui keranjang")
		}

		existingCart.Quantity = newQty
		return c.JSON(fiber.Map{
			"message": "Jumlah item diperbarui",
			"cart":    existingCart,
		})
	}

	newCart := models.Cart{
		CustomerID:       userUID,
		ProductVariantID: variant.ID,
		Quantity:         req.Quantity,
		IsSelected:       false,
		IsConverted:      false,
	}
	if err := database.DB.Create(&newCart).Error; err != nil {
		log.Println("[ERROR] Gagal create cart:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menambahkan ke keranjang")
	}

	return c.JSON(fiber.Map{
		"message": "Item ditambahkan ke keranjang",
		"cart":    newCart,
	})
}
