package handlers

import (
	"github.com/gofiber/fiber/v2"
	customer "github.com/winterheatherica/tokoaku-backend/internal/controllers/customer"
)

func CustomerRoutes(router fiber.Router) {
	router.Post("/cart", customer.AddToCart)
	router.Get("/cart", customer.GetCart)
}
