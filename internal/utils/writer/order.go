package writer

import (
	"context"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func CreateOrder(ctx context.Context, order *models.Order) error {
	return database.DB.WithContext(ctx).Create(order).Error
}

func CreateOrderLog(ctx context.Context, orderID uint) error {
	log := models.OrderLog{
		OrderID:  orderID,
		StatusID: 11,
	}
	return database.DB.WithContext(ctx).Create(&log).Error
}

func CreateOrderPromo(ctx context.Context, orderID uint, promoID uint, customerID string) error {
	if err := database.DB.WithContext(ctx).
		Create(&models.OrderPromo{OrderID: orderID, PromoID: promoID}).Error; err != nil {
		return err
	}
	return database.DB.WithContext(ctx).
		Model(&models.UserPromo{}).
		Where("promo_id = ? AND customer_id = ?", promoID, customerID).
		Update("redeemed", true).Error
}
