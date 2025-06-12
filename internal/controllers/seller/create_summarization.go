package seller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func callSummarizerFlask(reviews []string) (string, error) {
	baseURL := os.Getenv("MACHINE_LEARNING_BASE_URL")
	endpoint := baseURL + "/predict-summarize"

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

func CreatePositiveSummarization(c *fiber.Ctx) error {
	const sentimentPositiveID = 1

	productID := c.Params("id")

	var product models.Product
	if err := database.DB.Preload("Variants").Where("id = ?", productID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Produk tidak ditemukan",
		})
	}

	var variantIDs []string
	for _, v := range product.Variants {
		variantIDs = append(variantIDs, v.ID)
	}

	var reviews []models.Review
	if err := database.DB.
		Where("product_variant_id IN ?", variantIDs).
		Where("sentiment_id = ?", sentimentPositiveID).
		Order("created_at DESC").
		Limit(15).
		Find(&reviews).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil review",
		})
	}

	if len(reviews) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Belum ada cukup review positif untuk diringkas",
		})
	}

	var texts []string
	for _, r := range reviews {
		texts = append(texts, r.Text)
	}

	summaryText, err := callSummarizerFlask(texts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memanggil summarizer: " + err.Error(),
		})
	}

	var summary models.Summarization
	if err := database.DB.Where("product_id = ? AND sentiment_id = ?", product.ID, sentimentPositiveID).First(&summary).Error; err != nil {
		var count int64
		database.DB.Model(&models.Summarization{}).Where("product_id = ?", product.ID).Count(&count)
		if count >= 2 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Sudah mencapai batas 2 summarization untuk produk ini",
			})
		}

		summary = models.Summarization{
			ProductID:   product.ID,
			SentimentID: sentimentPositiveID,
			ReviewCount: uint(len(reviews)),
		}
		if err := database.DB.Create(&summary).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal menyimpan summarization",
			})
		}
	}

	detail := models.SummarizationDetail{
		SummarizationID: summary.ID,
		Text:            summaryText,
	}

	if err := database.DB.Create(&detail).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan hasil ringkasan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ringkasan review positif berhasil dibuat",
		"summary": summaryText,
	})
}
