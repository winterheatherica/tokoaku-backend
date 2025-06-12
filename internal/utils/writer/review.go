package writer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func SaveReview(review *models.Review) error {
	if err := database.DB.Create(&review).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save review to database")
	}
	return nil
}
