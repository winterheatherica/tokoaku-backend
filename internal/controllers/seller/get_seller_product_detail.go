package seller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetSellerProductDetail(c *fiber.Ctx) error {
	productID := c.Params("id")
	sellerID := c.Locals("uid").(string)
	ctx := context.Background()

	var product models.Product
	if err := database.DB.
		Where("id = ? AND seller_id = ?", productID, sellerID).
		First(&product).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	productForms, _ := fetcher.GetAllProductForms(ctx)
	productTypes, _ := fetcher.GetAllProductTypes(ctx)

	formMap := make(map[uint]string)
	for _, f := range productForms {
		formMap[f.ID] = f.Form
	}
	typeMap := make(map[uint]string)
	for _, t := range productTypes {
		typeMap[t.ID] = t.Name
	}

	response := map[string]interface{}{
		"id":                product.ID,
		"name":              product.Name,
		"description":       product.Description,
		"image_cover_url":   product.ImageCoverURL,
		"product_type_name": typeMap[product.ProductTypeID],
		"product_form_name": formMap[product.ProductFormID],
		"created_at":        product.CreatedAt,
		"updated_at":        product.UpdatedAt,
	}

	return c.JSON(response)
}
