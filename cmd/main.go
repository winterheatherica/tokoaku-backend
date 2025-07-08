package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/winterheatherica/tokoaku-backend/config"
	"github.com/winterheatherica/tokoaku-backend/internal/middleware"
	"github.com/winterheatherica/tokoaku-backend/internal/routes"
	"github.com/winterheatherica/tokoaku-backend/internal/services"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	} else {
		log.Println("[ENV]: ✅ .env file loaded")
	}

}

func main() {

	config.LoadAll()
	services.InitAll()
	// seed.RunAllSeeders(database.DB)
	// persistent.StartPersistentCacheRefresher()
	// volatile.StartVolatileCacheRefresher()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(middleware.Cors())

	routes.SetupRoutes(app)

	log.Printf("[MAIN]: 🚀 Server running on %s\n", config.App.BackendBaseURL)
	log.Fatal(app.Listen(":" + config.App.Port))
}
