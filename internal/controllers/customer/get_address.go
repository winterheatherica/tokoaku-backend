package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetActiveAddress(c *fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	var address models.Address
	if err := database.DB.
		Where("user_id = ? AND is_active = true", userID).
		First(&address).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Alamat aktif tidak ditemukan")
	}

	return c.JSON(fiber.Map{
		"message": "Alamat aktif ditemukan",
		"address": address,
	})
}
