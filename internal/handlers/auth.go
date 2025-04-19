package handlers

import (
	"github.com/gofiber/fiber/v2"
	auth "github.com/winterheatherica/tokoaku-backend/internal/controllers/auth"
)

func PublicAuthRoutes(router fiber.Router) {
	router.Post("/register", auth.Register)
	router.Post("/verify", auth.VerifyToken)
	router.Post("/login", auth.Login)
}

func PrivateAuthRoutes(router fiber.Router) {
	router.Get("/me", auth.Me)
}
