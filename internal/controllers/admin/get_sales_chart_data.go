package admin

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type SalesChartResponse struct {
	Historical []models.SalesData     `json:"historical"`
	Forecast   []models.SalesForecast `json:"forecast"`
	Analysis   string                 `json:"analysis"`
}

func GetSalesChartData(c *fiber.Ctx) error {
	db := database.DB
	ctx := context.Background()

	var forecastsDesc []models.SalesForecast
	if err := db.WithContext(ctx).
		Order("date DESC").
		Limit(30).
		Find(&forecastsDesc).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data forecast",
		})
	}

	if len(forecastsDesc) == 0 {
		return c.Status(http.StatusOK).JSON(SalesChartResponse{
			Historical: []models.SalesData{},
			Forecast:   []models.SalesForecast{},
			Analysis:   "",
		})
	}

	batchCount := map[string]int{}
	for _, f := range forecastsDesc {
		batchCount[f.BatchID]++
	}

	majorityBatch := forecastsDesc[0].BatchID
	maxCount := 0
	for batchID, count := range batchCount {
		if count > maxCount {
			majorityBatch = batchID
			maxCount = count
		}
	}

	var forecasts []models.SalesForecast
	if err := db.WithContext(ctx).
		Where("batch_id = ?", majorityBatch).
		Order("date ASC").
		Find(&forecasts).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data forecast berdasarkan batch",
		})
	}

	var historical []models.SalesData
	if err := db.WithContext(ctx).
		Order("date ASC").
		Find(&historical).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data historis",
		})
	}

	var batch models.SalesForecastBatch
	if err := db.WithContext(ctx).
		Where("id = ?", majorityBatch).
		First(&batch).Error; err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil analisis batch",
		})
	}

	return c.Status(http.StatusOK).JSON(SalesChartResponse{
		Historical: historical,
		Forecast:   forecasts,
		Analysis:   batch.Analysis,
	})
}
