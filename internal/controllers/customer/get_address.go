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

func GetAllAddress(c *fiber.Ctx) error {
	userID := c.Locals("uid").(string)

	var addresses []models.Address
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&addresses).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data alamat")
	}

	return c.JSON(fiber.Map{
		"message":   "Daftar alamat berhasil diambil",
		"addresses": addresses,
	})
}

func SetActiveAddress(c *fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	addressID := c.Params("id")

	if err := database.DB.Model(&models.Address{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal reset alamat aktif")
	}

	if err := database.DB.Model(&models.Address{}).
		Where("user_id = ? AND id = ?", userID, addressID).
		Update("is_active", true).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menetapkan alamat aktif")
	}

	return c.JSON(fiber.Map{
		"message": "Alamat aktif berhasil diperbarui",
	})
}
