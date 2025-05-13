package seller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func AddProductVariant(c *fiber.Ctx) error {
	productID := c.Params("id")

	type RequestBody struct {
		VariantName string `json:"variant_name"`
		Stock       uint   `json:"stock"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if body.VariantName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Variant name cannot be empty")
	}

	slug := utils.SlugifyText(body.VariantName)

	variant := models.ProductVariant{
		ID:          uuid.NewString(),
		VariantName: body.VariantName,
		ProductID:   productID,
		Stock:       body.Stock,
		Slug:        slug,
	}

	if err := database.DB.Create(&variant).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create product variant")
	}

	return c.JSON(fiber.Map{
		"message": "Product variant created successfully",
		"variant": variant,
	})
}
