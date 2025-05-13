package seller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func SetProductVariantPrice(c *fiber.Ctx) error {
	variantID := c.Params("id")
	if variantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Variant ID is required",
		})
	}

	var payload struct {
		Price uint `json:"price"`
	}
	if err := c.BodyParser(&payload); err != nil || payload.Price == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or missing price",
		})
	}

	var variant models.ProductVariant
	if err := database.DB.First(&variant, "id = ?", variantID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product variant not found",
		})
	}

	price := models.ProductPrice{
		ProductVariantID: variantID,
		Price:            payload.Price,
		CreatedAt:        time.Now(),
	}

	if err := database.DB.Create(&price).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store price",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Price added successfully",
		"price":   price,
	})
}
