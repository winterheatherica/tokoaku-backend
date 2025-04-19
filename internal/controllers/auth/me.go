package auth

import (
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

func Me(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*auth.Token)
	if !ok {
		log.Println("Tidak bisa ambil data user dari Firebase token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"uid":   token.UID,
			"email": token.Claims["email"],
		},
	})
}
