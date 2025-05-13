package auth

import (
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func Me(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*auth.Token)
	if !ok || token == nil {
		log.Println("Token Firebase tidak valid atau tidak ditemukan")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	uid := token.UID

	var user models.User
	if err := database.DB.First(&user, "id = ?", uid).Error; err != nil {
		log.Printf("User tidak ditemukan di database: %v\n", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan",
		})
	}

	roleName, err := fetcher.GetRoleNameByID(int(user.RoleID))
	if err != nil {
		log.Println("Gagal ambil nama role:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil role user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"phone":     user.Phone,
		"provider":  user.ProviderID,
		"role":      roleName,
		"name":      user.Name,
		"createdAt": user.CreatedAt,
	})
}
