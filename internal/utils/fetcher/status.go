package fetcher

import (
	"context"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetAllStatuses(ctx context.Context) ([]models.Status, error) {
	var statuses []models.Status
	err := database.DB.WithContext(ctx).
		Order("created_at ASC").
		Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
