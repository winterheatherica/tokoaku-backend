package handlers

import (
	"github.com/gofiber/fiber/v2"
	sellercontroller "github.com/winterheatherica/tokoaku-backend/internal/controllers/seller"
)

func SellerRoutes(router fiber.Router) {
	router.Get("/product-types", sellercontroller.GetAllProductTypes)
	router.Post("/products", sellercontroller.AddProduct)
}
