package seller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetAllBankAccounts(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	var accounts []models.BankAccount
	if err := database.DB.
		Preload("Bank").
		Where("user_id = ?", uid).
		Order("created_at DESC").
		Find(&accounts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data akun bank")
	}

	return c.JSON(fiber.Map{
		"message":  "Daftar akun bank berhasil diambil",
		"accounts": accounts,
	})
}
