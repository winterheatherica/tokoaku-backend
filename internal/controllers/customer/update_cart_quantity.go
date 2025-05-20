package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type UpdateCartQuantityRequest struct {
	ProductVariantID string `json:"product_variant_id"`
	Quantity         uint   `json:"quantity"`
}

func UpdateCartQuantity(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)

	var req UpdateCartQuantityRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format request tidak valid")
	}

	if req.ProductVariantID == "" || req.Quantity == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "product_variant_id dan quantity wajib diisi")
	}

	if err := database.DB.
		Model(&models.Cart{}).
		Where("customer_id = ? AND product_variant_id = ? AND is_converted = false", userUID, req.ProductVariantID).
		Update("quantity", req.Quantity).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengubah jumlah item")
	}

	return c.JSON(fiber.Map{
		"message": "Jumlah item berhasil diperbarui",
	})
}
