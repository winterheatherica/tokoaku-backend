package customer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/writer"
	"gorm.io/gorm"
)

func AddReview(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)
	fmt.Println("✅ UID:", userUID)

	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		fmt.Println("❌ BodyParser error:", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	fmt.Println("✅ Parsed review:", review)

	review.CustomerID = userUID

	if review.Rating < 1 || review.Rating > 5 {
		fmt.Println("❌ Invalid rating:", review.Rating)
		return fiber.NewError(fiber.StatusBadRequest, "rating must be between 1 and 5")
	}

	fmt.Println("📦 Checking if user has purchased variant:", review.ProductVariantID)
	hasPurchased, err := fetcher.HasUserPurchasedVariant(userUID, review.ProductVariantID)
	if err != nil {
		fmt.Println("❌ Error checking purchase:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !hasPurchased {
		fmt.Println("❌ User has not purchased this product variant")
		return fiber.NewError(fiber.StatusForbidden, "you have not purchased this product variant")
	}

	fmt.Println("🧠 Sending to sentiment analysis...")
	sentimentID, err := fetcher.AnalyzeSentiment(review.Text)
	if err != nil {
		fmt.Println("❌ Sentiment analysis failed:", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	fmt.Println("✅ Sentiment ID:", sentimentID)

	review.SentimentID = sentimentID
	review.CreatedAt = time.Now()

	fmt.Println("💾 Saving review to DB...")
	if err := writer.SaveReview(&review); err != nil {
		fmt.Println("❌ Failed to save review:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save review")
	}

	fmt.Println("🔎 Fetching product variant info...")
	variant, err := fetcher.GetProductVariantByID(c.Context(), review.ProductVariantID)
	if err != nil {
		fmt.Println("❌ Failed to fetch product variant:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch product variant")
	}

	fmt.Println("📊 Updating/creating summarization record...")
	var summarization models.Summarization
	if err := database.DB.
		Where("product_id = ? AND sentiment_id = ?", variant.ProductID, *sentimentID).
		First(&summarization).Error; err != nil {

		fmt.Println("ℹ️ No existing summarization. Creating new...")
		summarization = models.Summarization{
			ProductID:   variant.ProductID,
			SentimentID: *sentimentID,
			ReviewCount: 1,
		}
		if err := database.DB.Create(&summarization).Error; err != nil {
			fmt.Println("❌ Failed to create summarization:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to create summarization")
		}
	} else {
		fmt.Println("✏️ Updating review count...")
		if err := database.DB.Model(&models.Summarization{}).
			Where("id = ?", summarization.ID).
			UpdateColumn("review_count", gorm.Expr("review_count + ?", 1)).Error; err != nil {
			fmt.Println("❌ Failed to atomically update review count:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update summarization")
		}
	}

	// if summarization.ReviewCount%20 == 0 {
	// 	fmt.Println("📈 Time to summarize after", summarization.ReviewCount, "reviews...")

	// 	var variantIDs []string
	// 	if err := database.DB.Model(&models.ProductVariant{}).
	// 		Where("product_id = ?", variant.ProductID).
	// 		Pluck("id", &variantIDs).Error; err != nil {
	// 		fmt.Println("❌ Failed to get variant IDs:", err)
	// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to get variant IDs")
	// 	}

	// 	var reviews []models.Review
	// 	if err := database.DB.
	// 		Where("product_variant_id IN ?", variantIDs).
	// 		Where("sentiment_id = ?", *sentimentID).
	// 		Order("created_at DESC").
	// 		Limit(20).
	// 		Find(&reviews).Error; err != nil {
	// 		fmt.Println("❌ Failed to fetch reviews for summarization:", err)
	// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch recent reviews")
	// 	}

	// 	var texts []string
	// 	for _, r := range reviews {
	// 		texts = append(texts, r.Text)
	// 	}
	// 	fmt.Println("📤 Sending", len(texts), "reviews to summarizer...")

	// 	summaryText, err := callSummarizerFlask(texts)
	// 	if err != nil {
	// 		fmt.Println("❌ Failed to summarize:", err)
	// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate summarization: "+err.Error())
	// 	}
	// 	fmt.Println("✅ Summary result:", summaryText)

	// 	newDetail := models.SummarizationDetail{
	// 		SummarizationID: summarization.ID,
	// 		Text:            summaryText,
	// 	}
	// 	if err := database.DB.Create(&newDetail).Error; err != nil {
	// 		fmt.Println("❌ Failed to save summarization detail:", err)
	// 		return fiber.NewError(fiber.StatusInternalServerError, "failed to save summarization detail")
	// 	}
	// }

	if summarization.ReviewCount > 19 {
		// Jalankan summarization di background
		go func(productID string, sentimentID uint) {
			fmt.Println("📈 Running background summarization...")

			var variantIDs []string
			if err := database.DB.Model(&models.ProductVariant{}).
				Where("product_id = ?", productID).
				Pluck("id", &variantIDs).Error; err != nil {
				fmt.Println("❌ (BG) Failed to get variant IDs:", err)
				return
			}

			var reviews []models.Review
			if err := database.DB.
				Where("product_variant_id IN ?", variantIDs).
				Where("sentiment_id = ?", sentimentID).
				Order("created_at DESC").
				Limit(20).
				Find(&reviews).Error; err != nil {
				fmt.Println("❌ (BG) Failed to fetch reviews for summarization:", err)
				return
			}

			var texts []string
			for _, r := range reviews {
				texts = append(texts, r.Text)
			}

			summaryText, err := callSummarizerFlask(texts)
			if err != nil {
				fmt.Println("❌ (BG) Failed to summarize:", err)
				return
			}

			var summarization models.Summarization
			if err := database.DB.
				Where("product_id = ? AND sentiment_id = ?", productID, sentimentID).
				First(&summarization).Error; err != nil {
				fmt.Println("❌ (BG) Failed to find summarization again:", err)
				return
			}

			newDetail := models.SummarizationDetail{
				SummarizationID: summarization.ID,
				Text:            summaryText,
			}
			if err := database.DB.Create(&newDetail).Error; err != nil {
				fmt.Println("❌ (BG) Failed to save summarization detail:", err)
				return
			}

			fmt.Println("✅ (BG) Summarization completed and saved.")
		}(variant.ProductID, *sentimentID)
	}

	fmt.Println("📦 Fetching full review for return...")
	fullReview, err := fetcher.GetFullReview(review.ID)
	if err != nil {
		fmt.Println("❌ Failed to fetch full review:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch saved review")
	}

	fmt.Println("✅ Review process complete.")
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
