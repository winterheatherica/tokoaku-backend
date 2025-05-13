package seller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetSellerProducts(c *fiber.Ctx) error {
	sellerID := c.Locals("uid").(string)
	ctx := context.Background()

	var products []models.Product
	if err := database.DB.
		Where("seller_id = ?", sellerID).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch seller products")
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

	type ProductCard struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Description     string `json:"description"`
		ImageCoverURL   string `json:"image_cover_url"`
		ProductTypeName string `json:"product_type_name"`
		ProductFormName string `json:"product_form_name"`
	}

	var response []ProductCard
	for _, p := range products {
		response = append(response, ProductCard{
			ID:              p.ID,
			Name:            p.Name,
			Description:     p.Description,
			ImageCoverURL:   p.ImageCoverURL,
			ProductTypeName: typeMap[p.ProductTypeID],
			ProductFormName: formMap[p.ProductFormID],
		})
	}

	return c.JSON(response)
}
