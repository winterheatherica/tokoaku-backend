package customer

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type SelectCartItemRequest struct {
	ProductVariantID string `json:"product_variant_id"`
	IsSelected       bool   `json:"is_selected"`
}

func SelectCartItem(c *fiber.Ctx) error {
	log.Println("🔥 Masuk SelectCartItem")
	userUID := c.Locals("uid").(string)

	var req SelectCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format request tidak valid")
	}

	if req.ProductVariantID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "product_variant_id wajib diisi")
	}

	if err := database.DB.
		Model(&models.Cart{}).
		Where("customer_id = ? AND product_variant_id = ? AND is_converted = false", userUID, req.ProductVariantID).
		Update("is_selected", req.IsSelected).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengubah status is_selected")
	}

	log.Printf("Try update is_selected: user=%s variant=%s → %v", userUID, req.ProductVariantID, req.IsSelected)

	return c.JSON(fiber.Map{
		"message": "Berhasil mengubah status item di keranjang",
	})
}
