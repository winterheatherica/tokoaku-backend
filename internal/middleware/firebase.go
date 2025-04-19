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
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header tidak valid",
			})
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		client, err := firebase.App.Auth(context.Background())
		if err != nil {
			log.Println("Gagal inisialisasi Auth client:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Println("Token tidak valid:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid atau sudah expired",
			})
		}

		c.Locals("uid", token.UID)
		return c.Next()
	}
}
