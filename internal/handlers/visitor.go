package handlers

import (
	"github.com/gofiber/fiber/v2"
	visitor "github.com/winterheatherica/tokoaku-backend/internal/controllers/visitor"
)

func VisitorRoutes(router fiber.Router) {
	router.Get("/products", visitor.GetVisitorProductList)
	router.Get("/products/:slug", visitor.GetVisitorProductDetail)
	router.Get("/product/:productSlug/variant/:variantSlug", visitor.GetVisitorProductDetailByVariant)
	router.Get("/product-reference", visitor.GetvisitorProductReferenceData)
	router.Get("/discount-banner", visitor.GetDiscountBanners)
	router.Get("/highlighted-product", visitor.GetHighlightedProductCards)
	router.Get("/products-by-form", visitor.GetProductByForm)
	router.Get("/products-by-type", visitor.GetProductByType)

	router.Get("/cloudinary-prefix", visitor.GetCloudinaryPublicImagePrefix)
	router.Get("/product-types", visitor.GetProductTypes)
	router.Get("/product-forms", visitor.GetProductForms)
	router.Get("/current-event", visitor.GetDiscountLimit)
	router.Get("/reviews/:product_slug", visitor.GetReviewsByProduct)
}
