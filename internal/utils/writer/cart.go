package writer

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func UpdateCartQuantity(existingCart *models.Cart, addQuantity uint) (*models.Cart, error) {
	newQty := existingCart.Quantity + addQuantity

	if err := database.DB.
		Model(&models.Cart{}).
		Where("customer_id = ? AND product_variant_id = ? AND is_converted = false",
			existingCart.CustomerID, existingCart.ProductVariantID).
		Update("quantity", newQty).Error; err != nil {
		log.Printf("[DB] Failed to update cart quantity: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update cart quantity")
	}

	existingCart.Quantity = newQty
	return existingCart, nil
}

func InsertNewCart(newCart *models.Cart) error {
	if err := database.DB.Create(newCart).Error; err != nil {
		log.Printf("[DB] Failed to insert new cart: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to add item to cart")
	}
	return nil
}

func MarkCartAsConverted(ctx context.Context, customerID string) error {
	return database.DB.WithContext(ctx).
		Model(&models.Cart{}).
		Where("customer_id = ? AND is_selected = true AND is_converted = false", customerID).
		Update("is_converted", true).Error
}
