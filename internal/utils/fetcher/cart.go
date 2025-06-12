package fetcher

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetUnconvertedCart(userUID string, productVariantID string) (*models.Cart, error) {
	var cart models.Cart
	if err := database.DB.
		Where("customer_id = ? AND product_variant_id = ? AND is_converted = false", userUID, productVariantID).
		First(&cart).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "no unconverted cart found for this variant")
	}
	return &cart, nil
}

func GetUnconvertedCarts(ctx context.Context, customerID string) ([]models.Cart, error) {
	var carts []models.Cart
	if err := database.DB.WithContext(ctx).
		Preload("ProductVariant.Product.Seller").
		Preload("ProductVariant.Product").
		Preload("ProductVariant").
		Where("customer_id = ? AND is_converted = false", customerID).
		Order("created_at DESC").
		Find(&carts).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve unconverted carts")
	}
	return carts, nil
}

func GetSelectedUnconvertedCarts(ctx context.Context, customerID string) ([]models.Cart, error) {
	var carts []models.Cart
	if err := database.DB.WithContext(ctx).
		Preload("ProductVariant.Product").
		Preload("ProductVariant").
		Where("customer_id = ? AND is_selected = true AND is_converted = false", customerID).
		Find(&carts).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "gagal mengambil cart")
	}
	return carts, nil
}
