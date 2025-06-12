package fetcher

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetOrderShippingsWithItemsByOrderID(ctx context.Context, orderID uint) ([]models.OrderShipping, error) {
	var shippings []models.OrderShipping

	err := database.DB.WithContext(ctx).
		Preload("ShippingOption").
		Preload("Seller").
		Preload("BankAccount").
		Preload("OrderItems.ProductVariant.Product").
		Where("order_id = ?", orderID).
		Find(&shippings).Error

	if err != nil {
		return nil, err
	}
	return shippings, nil
}

func HasUserPurchasedVariant(userUID string, variantID string) (bool, error) {
	var orders []models.Order
	if err := database.DB.Where("customer_id = ?", userUID).Find(&orders).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve orders")
	}
	if len(orders) == 0 {
		return false, nil
	}

	var paidOrderIDs []uint
	for _, o := range orders {
		var logs []models.OrderLog
		if err := database.DB.Where("order_id = ?", o.ID).Order("created_at DESC").Find(&logs).Error; err == nil && len(logs) > 0 {
			if logs[0].StatusID == 12 {
				paidOrderIDs = append(paidOrderIDs, o.ID)
			}
		}
	}
	if len(paidOrderIDs) == 0 {
		return false, nil
	}

	var shippings []models.OrderShipping
	if err := database.DB.Where("order_id IN ?", paidOrderIDs).Find(&shippings).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve order shipping data")
	}
	if len(shippings) == 0 {
		return false, nil
	}

	var shippingIDs []uint
	for _, s := range shippings {
		shippingIDs = append(shippingIDs, s.ID)
	}

	var selectedVariant models.ProductVariant
	if err := database.DB.Where("id = ?", variantID).First(&selectedVariant).Error; err != nil {
		return false, fiber.NewError(fiber.StatusBadRequest, "variant not found")
	}

	var sameProductVariantIDs []string
	if err := database.DB.
		Model(&models.ProductVariant{}).
		Where("product_id = ?", selectedVariant.ProductID).
		Pluck("id", &sameProductVariantIDs).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve product variant list")
	}

	var items []models.OrderItem
	if err := database.DB.
		Where("order_shipping_id IN ?", shippingIDs).
		Where("product_variant_id IN ?", sameProductVariantIDs).
		Find(&items).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve order items")
	}
	if len(items) == 0 {
		return false, nil
	}
	return true, nil
}
