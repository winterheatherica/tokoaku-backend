package auth

import (
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func UpdateMe(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*auth.Token)
	if !ok || token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var body struct {
		Username *string `json:"username"`
		Name     *string `json:"name"`
		Phone    *string `json:"phone"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data tidak valid",
		})
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", token.UID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan",
		})
	}

	if body.Username != nil {
		user.Username = body.Username
	}
	if body.Name != nil {
		user.Name = body.Name
	}
	if body.Phone != nil {
		user.Phone = body.Phone
	}

	if err := database.DB.Save(&user).Error; err != nil {
		log.Println("Gagal update user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan perubahan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil diperbarui",
		"user":    user,
	})
}
