package fetcher

import (
	"context"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetOrderLogsByOrderID(ctx context.Context, orderID uint) ([]models.OrderLog, error) {
	var logs []models.OrderLog
	err := database.DB.WithContext(ctx).
		Preload("Status").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}
