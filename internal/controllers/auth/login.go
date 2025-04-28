package auth

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	firebaseService "github.com/winterheatherica/tokoaku-backend/internal/services/firebase"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func Login(c *fiber.Ctx) error {
	log.Println("[BACKEND] 🚀 Masuk ke Login Handler")

	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		log.Println("[BACKEND] Gagal parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid",
		})
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		log.Println("[BACKEND] User tidak ditemukan:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Email atau password salah",
		})
	}

	roleName, err := utils.GetRoleNameByID(int(user.RoleID))
	if err != nil {
		log.Println("[BACKEND] Gagal ambil role:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil role user",
		})
	}

	ctx := context.Background()

	claims := map[string]interface{}{
		"role": roleName,
	}

	if err := firebaseService.FirebaseAuth.SetCustomUserClaims(ctx, user.ID, claims); err != nil {
		log.Println("[BACKEND] Gagal set custom claims:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyetel role user",
		})
	}

	token, err := firebaseService.FirebaseAuth.CustomToken(ctx, user.ID)
	if err != nil {
		log.Println("[BACKEND] Gagal membuat custom token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Login berhasil",
		"customToken": token,
		"email":       user.Email,
		"uid":         user.ID,
		"role":        roleName,
	})
}
