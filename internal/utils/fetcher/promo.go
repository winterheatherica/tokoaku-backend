package fetcher

import (
	"context"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetAllActivePromos(ctx context.Context) ([]models.Promo, error) {
	now := time.Now()
	var promos []models.Promo

	if err := database.DB.
		Preload("ValueType").
		Where("start_at <= ? AND end_at >= ?", now, now).
		Order("start_at ASC").
		Find(&promos).Error; err != nil {
		return nil, err
	}

	return promos, nil
}
