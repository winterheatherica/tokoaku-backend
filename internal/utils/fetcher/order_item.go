package fetcher

import (
	"context"

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
