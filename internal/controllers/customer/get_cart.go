package customer

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

type CartItemResponse struct {
	ProductVariantID     string    `json:"product_variant_id"`
	ProductName          string    `json:"product_name"`
	ProductSlug          string    `json:"product_slug"`
	ProductVariantName   string    `json:"product_variant_name"`
	ProductVariantSlug   string    `json:"product_variant_slug"`
	Quantity             uint      `json:"quantity"`
	ImageURL             string    `json:"image_url"`
	ProductVariantImages []string  `json:"product_variant_images"`
	AddedAt              time.Time `json:"added_at"`
	IsSelected           bool      `json:"is_selected"`
}

type ProductGroup struct {
	ProductName string             `json:"product_name"`
	CartItems   []CartItemResponse `json:"cart_items"`
}

type SellerGroup struct {
	SellerName string         `json:"seller_name"`
	Products   []ProductGroup `json:"products"`
}

func GetGroupedCart(c *fiber.Ctx) error {
	customerID := c.Locals("uid").(string)

	cartItems, err := fetcher.GetUnconvertedCarts(c.Context(), customerID)
	if err != nil {
		return err
	}

	grouped := make(map[string]map[string][]CartItemResponse)

	for _, item := range cartItems {
		variantImages, err := fetcher.GetAllVariantImages(c.Context(), item.ProductVariantID)
		if err != nil {
			variantImages = []models.ProductVariantImage{}
		}

		imageURL := item.ProductVariant.Product.ImageCoverURL
		var images []string
		for _, img := range variantImages {
			images = append(images, img.ImageURL)
			if img.IsVariantCover {
				imageURL = img.ImageURL
			}
		}

		sellerName := "-"
		if item.ProductVariant.Product.Seller.Name != nil {
			sellerName = *item.ProductVariant.Product.Seller.Name
		}
		productSlug := item.ProductVariant.Product.Slug

		cartData := CartItemResponse{
			ProductVariantID:     item.ProductVariantID,
			ProductName:          item.ProductVariant.Product.Name,
			ProductSlug:          productSlug,
			ProductVariantName:   item.ProductVariant.VariantName,
			ProductVariantSlug:   item.ProductVariant.Slug,
			Quantity:             item.Quantity,
			ImageURL:             imageURL,
			ProductVariantImages: images,
			AddedAt:              item.CreatedAt,
			IsSelected:           item.IsSelected,
		}

		if _, ok := grouped[sellerName]; !ok {
			grouped[sellerName] = make(map[string][]CartItemResponse)
		}
		grouped[sellerName][productSlug] = append(grouped[sellerName][productSlug], cartData)
	}

	var response []SellerGroup
	for sellerName, products := range grouped {
		var groupedProducts []ProductGroup
		for _, items := range products {
			productName := ""
			if len(items) > 0 {
				productName = items[0].ProductName
			}

			groupedProducts = append(groupedProducts, ProductGroup{
				ProductName: productName,
				CartItems:   items,
			})
		}

		response = append(response, SellerGroup{
			SellerName: sellerName,
			Products:   groupedProducts,
		})
	}

	return c.JSON(response)
}
