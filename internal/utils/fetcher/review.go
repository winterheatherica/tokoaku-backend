package fetcher

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func AnalyzeSentiment(text string) (*uint, error) {
	mlURL := os.Getenv("MACHINE_LEARNING_BASE_URL") + "/analyze-sentiment"
	payload := map[string]string{"text": text}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := http.Post(mlURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to connect to sentiment API")
	}
	defer resp.Body.Close()

	var result struct {
		Label      string  `json:"label"`
		Confidence float64 `json:"confidence"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "invalid response from sentiment API")
	}

	switch result.Label {
	case "Positive":
		id := uint(1)
		return &id, nil
	case "Negative":
		id := uint(2)
		return &id, nil
	default:
		return nil, nil
	}
}

func GetFullReview(reviewID uint) (*models.Review, error) {
	var fullReview models.Review
	if err := database.DB.
		Preload("Customer").
		Preload("ProductVariant").
		Preload("Sentiment").
		First(&fullReview, reviewID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve full review from database")
	}
	return &fullReview, nil
}
