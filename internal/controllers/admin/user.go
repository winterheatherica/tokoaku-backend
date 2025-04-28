package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetRecentUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := database.DB.
		Select("*").
		Omit("password_hash").
		Preload("Role").
		Preload("Provider").
		Order("created_at DESC").
		Limit(10).
		Find(&users).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data user",
		})
	}

	return c.JSON(users)
}
