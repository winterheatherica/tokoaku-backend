package fetcher

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetAddressByIDAndUserID(ctx context.Context, addressID uuid.UUID, userID string) (*models.Address, error) {
	var address models.Address
	if err := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", addressID, userID).
		First(&address).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "alamat tidak ditemukan")
	}
	return &address, nil
}
