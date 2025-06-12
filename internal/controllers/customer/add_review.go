package customer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/writer"
)

func AddReview(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)

	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	review.CustomerID = userUID

	if review.Rating < 1 || review.Rating > 5 {
		return fiber.NewError(fiber.StatusBadRequest, "rating must be between 1 and 5")
	}

	hasPurchased, err := fetcher.HasUserPurchasedVariant(userUID, review.ProductVariantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !hasPurchased {
		return fiber.NewError(fiber.StatusForbidden, "you have not purchased this product variant")
	}

	sentimentID, err := fetcher.AnalyzeSentiment(review.Text)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	review.SentimentID = sentimentID
	review.CreatedAt = time.Now()

	if err := writer.SaveReview(&review); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save review")
	}

	variant, err := fetcher.GetProductVariantByID(c.Context(), review.ProductVariantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch product variant")
	}

	var summarization models.Summarization
	if err := database.DB.
		Where("product_id = ? AND sentiment_id = ?", variant.ProductID, *sentimentID).
		First(&summarization).Error; err != nil {

		summarization = models.Summarization{
			ProductID:   variant.ProductID,
			SentimentID: *sentimentID,
			ReviewCount: 1,
		}
		if err := database.DB.Create(&summarization).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create summarization")
		}
	} else {
		summarization.ReviewCount += 1
		if err := database.DB.Save(&summarization).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update summarization")
		}
	}

	if summarization.ReviewCount%15 == 0 {
		var variantIDs []string
		if err := database.DB.Model(&models.ProductVariant{}).
			Where("product_id = ?", variant.ProductID).
			Pluck("id", &variantIDs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to get variant IDs")
		}

		var reviews []models.Review
		if err := database.DB.
			Where("product_variant_id IN ?", variantIDs).
			Where("sentiment_id = ?", *sentimentID).
			Order("created_at DESC").
			Limit(15).
			Find(&reviews).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch recent reviews for summarization")
		}

		var texts []string
		for _, r := range reviews {
			texts = append(texts, r.Text)
		}

		summaryText, err := callSummarizerFlask(texts)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to generate summarization: "+err.Error())
		}

		newDetail := models.SummarizationDetail{
			SummarizationID: summarization.ID,
			Text:            summaryText,
		}
		if err := database.DB.Create(&newDetail).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to save summarization detail")
		}
	}

	fullReview, err := fetcher.GetFullReview(review.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch saved review")
	}

	return c.JSON(fiber.Map{
		"message": "review submitted successfully",
		"review":  fullReview,
	})
}

func callSummarizerFlask(reviews []string) (string, error) {
	baseURL := os.Getenv("MACHINE_LEARNING_BASE_URL")
	endpoint := baseURL + "/create-summarize"

	body := map[string]interface{}{"reviews": reviews}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var result map[string]string
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	return result["summary"], nil
}
