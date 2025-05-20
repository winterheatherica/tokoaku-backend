package fetcher

import (
	"context"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetShippingStatusesByOrderShippingID(ctx context.Context, orderShippingID uint) ([]models.OrderShippingStatus, error) {
	var statuses []models.OrderShippingStatus
	err := database.DB.WithContext(ctx).
		Preload("Status").
		Where("order_shipping_id = ?", orderShippingID).
		Order("created_at ASC").
		Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
