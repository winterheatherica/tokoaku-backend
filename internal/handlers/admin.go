package handlers

import (
	"github.com/gofiber/fiber/v2"
	admin "github.com/winterheatherica/tokoaku-backend/internal/controllers/admin"
)

func AdminRoutes(router fiber.Router) {
	router.Post("/discount", admin.AddDiscount)
	router.Post("/upload-sales", admin.UploadSalesCSV)
	router.Post("/predict-sales", admin.PredictSales)
	router.Get("/sales-forecast-history", admin.GetSalesChartData)
	router.Get("/users", admin.GetRecentUsers)
}
