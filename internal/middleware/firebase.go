package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/services/firebase"
)

func VerifyFirebaseToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			log.Println("[ERROR] Authorization header kosong")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header kosong",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("[ERROR] Authorization header tidak diawali Bearer")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header tidak valid",
			})
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		idToken = strings.TrimSpace(idToken)

		if idToken == "" {
			log.Println("[ERROR] ID Token kosong setelah Bearer")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak ditemukan",
			})
		}

		client, err := firebase.App.Auth(context.Background())
		if err != nil {
			log.Println("[ERROR] Gagal inisialisasi Firebase Auth client:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menginisialisasi Auth client",
			})
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Println("[ERROR] Token tidak valid:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		c.Locals("uid", token.UID)
		// c.Locals("user", token)

		return c.Next()
	}
}
