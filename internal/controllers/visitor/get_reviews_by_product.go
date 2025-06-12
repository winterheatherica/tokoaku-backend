package visitor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetReviewsByProduct(c *fiber.Ctx) error {
	productSlug := c.Params("product_slug")

	var product models.Product
	if err := database.DB.
		Where("slug = ?", productSlug).
		First(&product).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Produk tidak ditemukan")
	}

	var variants []models.ProductVariant
	if err := database.DB.
		Where("product_id = ?", product.ID).
		Find(&variants).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil varian produk")
	}

	if len(variants) == 0 {
		return c.JSON([]models.Review{})
	}

	var variantIDs []string
	for _, v := range variants {
		variantIDs = append(variantIDs, v.ID)
	}

	var reviews []models.Review
	if err := database.DB.
		Preload("Customer").
		Preload("ProductVariant").
		Preload("Sentiment").
		Where("product_variant_id IN ?", variantIDs).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil review")
	}

	var response []fiber.Map
	for _, r := range reviews {
		response = append(response, fiber.Map{
			"id":           r.ID,
			"text":         r.Text,
			"rating":       r.Rating,
			"created_at":   r.CreatedAt,
			"label":        r.Sentiment.Name,
			"variant_name": r.ProductVariant.VariantName,
			"customer": fiber.Map{
				"name": r.Customer.Username,
			},
		})
	}

	return c.JSON(response)

}
