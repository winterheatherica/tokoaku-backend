package fetcher

import (
	"context"
	"errors"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetUserPromoByPromoID(ctx context.Context, userID string, promoID uint) (*models.UserPromo, error) {
	var userPromo models.UserPromo

	if err := database.DB.
		Preload("Promo.ValueType").
		Where("customer_id = ? AND promo_id = ?", userID, promoID).
		First(&userPromo).Error; err != nil {
		return nil, err
	}

	return &userPromo, nil
}

func GetAvailablePromosForUser(ctx context.Context, userID string) ([]models.Promo, error) {
	now := time.Now()
	var promos []models.Promo

	subQuery := database.DB.
		Model(&models.UserPromo{}).
		Select("promo_id").
		Where("customer_id = ?", userID)

	if err := database.DB.
		Preload("ValueType").
		Where("id NOT IN (?) AND start_at <= ? AND end_at >= ?", subQuery, now, now).
		Find(&promos).Error; err != nil {
		return nil, err
	}

	return promos, nil
}

func CalculatePromoDiscount(promo models.Promo, subtotalOriginalPrice uint) (uint, error) {
	if subtotalOriginalPrice < promo.MinPriceValue {
		return 0, errors.New("subtotal tidak memenuhi minimum harga promo")
	}

	if promo.ValueType.Name == "" {
		return 0, errors.New("tipe value promo tidak ditemukan")
	}

	valueType := promo.ValueType.Name

	switch valueType {
	case "Percentage", "percentage", "percent", "persen":
		discount := (subtotalOriginalPrice * promo.Value) / 100

		if promo.MaxValue > 0 && discount > promo.MaxValue {
			discount = promo.MaxValue
		}
		return discount, nil

	case "flat", "nominal":
		return promo.Value, nil

	default:
		return 0, errors.New("tipe promo tidak dikenali: " + valueType)
	}
}
