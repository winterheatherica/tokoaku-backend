package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/winterheatherica/tokoaku-backend/config"
)

func Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: config.App.FrontendBaseURL,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})
}
