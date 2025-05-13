package handlers

import (
	"github.com/gofiber/fiber/v2"
	admin "github.com/winterheatherica/tokoaku-backend/internal/controllers/admin"
)

func AdminRoutes(router fiber.Router) {
	router.Get("/users", admin.GetRecentUsers)
}
