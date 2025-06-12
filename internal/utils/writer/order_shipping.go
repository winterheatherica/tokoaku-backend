package writer

import (
	"context"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func CreateOrderShippingAndItems(
	ctx context.Context,
	orderID uint,
	orderShippings []models.OrderShipping,
	groupedCarts map[string][]models.Cart,
) error {
	for _, shipping := range orderShippings {
		if err := database.DB.WithContext(ctx).Create(&shipping).Error; err != nil {
			continue
		}

		_ = database.DB.WithContext(ctx).Create(&models.OrderShippingStatus{
			OrderShippingID: shipping.ID,
			StatusID:        21,
		}).Error

		items := groupedCarts[shipping.SellerID]
		for _, item := range items {
			_ = database.DB.WithContext(ctx).Create(&models.OrderItem{
				OrderShippingID:  shipping.ID,
				ProductVariantID: item.ProductVariantID,
				Quantity:         item.Quantity,
			}).Error
		}
	}
	return nil
}
