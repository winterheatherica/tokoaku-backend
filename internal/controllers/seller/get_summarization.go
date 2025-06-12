package seller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetPositiveSummarization(c *fiber.Ctx) error {
	const sentimentPositiveID = 1

	productID := c.Params("id")

	var product models.Product
	if err := database.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Produk tidak ditemukan",
		})
	}

	var summary models.Summarization
	if err := database.DB.
		Where("product_id = ? AND sentiment_id = ?", product.ID, sentimentPositiveID).
		First(&summary).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ringkasan review positif belum tersedia",
		})
	}

	var latestDetail models.SummarizationDetail
	if err := database.DB.
		Where("summarization_id = ?", summary.ID).
		Order("created_at DESC").
		First(&latestDetail).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Detail ringkasan tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"product_id":     summary.ProductID,
		"sentiment_id":   summary.SentimentID,
		"review_count":   summary.ReviewCount,
		"latest_summary": latestDetail.Text,
	})
}

func GetNegativeSummarization(c *fiber.Ctx) error {
	const sentimentNegativeID = 2

	productID := c.Params("id")

	var product models.Product
	if err := database.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Produk tidak ditemukan",
		})
	}

	var summary models.Summarization
	if err := database.DB.
		Where("product_id = ? AND sentiment_id = ?", product.ID, sentimentNegativeID).
		First(&summary).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Ringkasan review negatif belum tersedia",
		})
	}

	var latestDetail models.SummarizationDetail
	if err := database.DB.
		Where("summarization_id = ?", summary.ID).
		Order("created_at DESC").
		First(&latestDetail).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Detail ringkasan tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"product_id":     summary.ProductID,
		"sentiment_id":   summary.SentimentID,
		"review_count":   summary.ReviewCount,
		"latest_summary": latestDetail.Text,
	})
}
