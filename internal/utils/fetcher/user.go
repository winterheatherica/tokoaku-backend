package fetcher

import (
	"context"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetUserByUID(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	if err := database.DB.WithContext(ctx).Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
