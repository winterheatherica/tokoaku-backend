package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/writer"
)

type AddToCartInput struct {
	ProductVariantID string `json:"product_variant_id"`
	Quantity         uint   `json:"quantity"`
}

func AddToCart(c *fiber.Ctx) error {
	customerID := c.Locals("uid").(string)

	var input AddToCartInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if input.Quantity == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "quantity must be at least 1")
	}

	productVariant, err := fetcher.GetProductVariantByID(c.Context(), input.ProductVariantID)
	if err != nil {
		return err
	}

	existingCart, err := fetcher.GetUnconvertedCart(customerID, productVariant.ID)
	if err == nil {
		updatedCart, err := writer.UpdateCartQuantity(existingCart, input.Quantity)
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"message": "cart quantity updated",
			"cart":    updatedCart,
		})
	}

	newCart := &models.Cart{
		CustomerID:       customerID,
		ProductVariantID: productVariant.ID,
		Quantity:         input.Quantity,
		IsSelected:       false,
		IsConverted:      false,
	}

	if err := writer.InsertNewCart(newCart); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "item added to cart",
		"cart":    newCart,
	})
}
