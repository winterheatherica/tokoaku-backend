package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/handlers"
	"github.com/winterheatherica/tokoaku-backend/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	public := app.Group("/auth")
	handlers.PublicAuthRoutes(public)

	protected := app.Group("/auth", middleware.VerifyFirebaseToken())
	handlers.PrivateAuthRoutes(protected)
}
