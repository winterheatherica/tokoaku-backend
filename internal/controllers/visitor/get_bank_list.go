package visitor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetBankList(c *fiber.Ctx) error {
	var bankList []models.BankList

	if err := database.DB.Order("code ASC").Find(&bankList).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data bank")
	}

	return c.JSON(fiber.Map{
		"message": "Data bank berhasil diambil",
		"data":    bankList,
	})
}
