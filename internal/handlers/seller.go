package handlers

import (
	"github.com/gofiber/fiber/v2"
	seller "github.com/winterheatherica/tokoaku-backend/internal/controllers/seller"
)

func SellerRoutes(router fiber.Router) {
	router.Post("/products", seller.AddProduct)
	router.Post("/products/upload-image", seller.UploadProductImage)
	router.Post("/products/:id/variants", seller.AddProductVariant)
	router.Post("/variants/:id/images", seller.UploadProductVariantImage)
	router.Post("/variants/:id/price", seller.SetProductVariantPrice)
	router.Post("/variants/:id/cover", seller.SetVariantCover)
	router.Post("/products/:id/summarize/positive", seller.CreatePositiveSummarization)
	router.Post("/bank-account/add", seller.AddBankAccount)

	router.Get("/products", seller.GetSellerProducts)
	router.Get("/products/:id", seller.GetSellerProductDetail)
	router.Get("/products/:id/variants", seller.GetProductVariants)
	router.Get("/variants/:id/images", seller.GetProductVariantImages)
	router.Get("/variants/:id/price", seller.GetLatestProductVariantPrice)
	router.Get("/products/:id/summarize/positive", seller.GetPositiveSummarization)
	router.Get("/products/:id/summarize/negative", seller.GetNegativeSummarization)
	router.Get("/bank-account/all", seller.GetAllBankAccounts)

	router.Get("/sales-forecast-history", seller.GetSalesChartData)

	router.Patch("/bank-account/set-active/:id", seller.SetActiveBankAccount)

}
