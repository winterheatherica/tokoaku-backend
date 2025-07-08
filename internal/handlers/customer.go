package handlers

import (
	"github.com/gofiber/fiber/v2"
	customer "github.com/winterheatherica/tokoaku-backend/internal/controllers/customer"
)

func CustomerRoutes(router fiber.Router) {
	router.Post("/cart", customer.AddToCart)
	router.Post("/order/create", customer.CreateOrder)
	router.Post("/purchase/:order_id", customer.PurchaseOrderDemo)
	router.Post("/review", customer.AddReview)
	router.Post("/address/add", customer.AddAddress)

	router.Get("/review/check/:product_slug", customer.LetReview)
	router.Get("/cart/grouped", customer.GetGroupedCart)
	router.Get("/cart/fees", customer.GetSellerFee)
	router.Get("/checkout/preview", customer.PreviewCheckout)
	router.Get("/address", customer.GetActiveAddress)
	router.Get("/order/:id", customer.GetOrderDetail)
	router.Get("/address/all", customer.GetAllAddress)

	router.Patch("/cart/select", customer.SelectCartItem)
	router.Patch("/cart/quantity", customer.UpdateCartQuantity)
	router.Patch("/address/set-active/:id", customer.SetActiveAddress)

}
