package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/controllers/auth"
	"github.com/winterheatherica/tokoaku-backend/internal/handlers"
	"github.com/winterheatherica/tokoaku-backend/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	public := app.Group("/auth")
	handlers.PublicAuthRoutes(public)

	protected := app.Group("/auth", middleware.VerifyFirebaseToken())
	handlers.PrivateAuthRoutes(protected)
	protected.Patch("/me", auth.UpdateMe)

	admin := app.Group("/admin", middleware.VerifyFirebaseToken())
	handlers.AdminRoutes(admin)

	seller := app.Group("/seller", middleware.VerifyFirebaseToken()) // Pastikan 'seller' group ini digunakan
	handlers.SellerRoutes(seller)

}
