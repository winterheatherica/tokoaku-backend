package auth

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	firebaseService "github.com/winterheatherica/tokoaku-backend/internal/services/firebase"
)

func Login(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau password salah"})
	}

	ctx := context.Background()

	claims := map[string]interface{}{
		"role": user.Role,
	}

	err := firebaseService.FirebaseAuth.SetCustomUserClaims(ctx, user.ID, claims)
	if err != nil {
		log.Println("Gagal set custom claims:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyetel role user",
		})
	}

	token, err := firebaseService.FirebaseAuth.CustomToken(ctx, user.ID)
	if err != nil {
		log.Println("Gagal buat custom token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuat token"})
	}

	log.Println("Login berhasil untuk:", body.Email)

	return c.JSON(fiber.Map{
		"message":     "Login berhasil",
		"customToken": token,
		"email":       user.Email,
		"uid":         user.ID,
	})
}
