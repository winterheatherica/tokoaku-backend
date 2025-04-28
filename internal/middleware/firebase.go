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
		// Ambil Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Println("[ERROR] Authorization header kosong")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header kosong",
			})
		}

		// Pastikan header diawali "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("[ERROR] Authorization header tidak diawali Bearer")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header tidak valid",
			})
		}

		// Ambil ID Token
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		idToken = strings.TrimSpace(idToken) // trims extra spaces

		if idToken == "" {
			log.Println("[ERROR] ID Token kosong setelah Bearer")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak ditemukan",
			})
		}

		// Ambil client Firebase
		client, err := firebase.App.Auth(context.Background())
		if err != nil {
			log.Println("[ERROR] Gagal inisialisasi Firebase Auth client:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menginisialisasi Auth client",
			})
		}

		// Verifikasi Token
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Println("[ERROR] Token tidak valid:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// Token valid, simpan data user di context
		c.Locals("uid", token.UID)
		c.Locals("user", token)

		log.Println("[INFO] Token valid. UID:", token.UID)

		// Lanjutkan ke handler berikutnya
		return c.Next()
	}
}
