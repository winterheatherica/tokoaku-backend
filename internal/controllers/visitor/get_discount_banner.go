package visitor

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetDiscountBanners(c *fiber.Ctx) error {
	db := database.DB
	now := time.Now()

	var discounts []models.Discount
	if err := db.
		Preload("ValueType").
		Preload("DiscountSponsor").
		Preload("CloudService").
		Where("start_at <= ? AND end_at >= ?", now, now).
		Order("start_at DESC").
		Find(&discounts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch active discounts: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"data":    discounts,
	})
}
