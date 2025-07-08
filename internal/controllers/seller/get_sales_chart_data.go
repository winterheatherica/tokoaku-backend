package seller

import (
	"context"
	"fmt"
	"net/http"
	"sort"

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

	// Step 1: Ambil 30 forecast terbaru tanpa filter batch
	var forecastsDesc []models.SalesForecast
	if err := db.WithContext(ctx).
		Order("date DESC").
		Limit(30).
		Find(&forecastsDesc).Error; err != nil {
		fmt.Println("❌ Error ambil 30 forecast desc:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data forecast",
		})
	}
	fmt.Printf("✅ Jumlah forecast diambil (desc): %d\n", len(forecastsDesc))

	if len(forecastsDesc) == 0 {
		fmt.Println("⚠️ Tidak ada data forecast.")
		return c.Status(http.StatusOK).JSON(SalesChartResponse{
			Historical: []models.SalesData{},
			Forecast:   []models.SalesForecast{},
			Analysis:   "",
		})
	}

	// Step 2: Reverse jadi ASC
	sort.Slice(forecastsDesc, func(i, j int) bool {
		return forecastsDesc[i].Date.Before(forecastsDesc[j].Date)
	})
	forecasts := forecastsDesc

	// Step 3: Ambil tanggal paling awal dan akhir
	earliestForecastDate := forecasts[0].Date
	latestForecastDate := forecasts[len(forecasts)-1].Date
	fmt.Printf("📅 Tanggal forecast paling awal: %s\n", earliestForecastDate.Format("2006-01-02"))
	fmt.Printf("📅 Tanggal forecast paling akhir: %s\n", latestForecastDate.Format("2006-01-02"))

	// Step 4: Ambil historical sebelum tanggal awal forecast
	var historical []models.SalesData
	if err := db.WithContext(ctx).
		Where("date < ?", earliestForecastDate).
		Order("date ASC").
		Find(&historical).Error; err != nil {
		fmt.Println("❌ Error ambil historical:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data historis",
		})
	}
	fmt.Printf("✅ Jumlah historical data: %d (sampai sebelum %s)\n", len(historical), earliestForecastDate.Format("2006-01-02"))

	// Step 5: Ambil analisis dari batch pertama dalam data forecast
	batchID := forecasts[0].BatchID
	var batch models.SalesForecastBatch
	if err := db.WithContext(ctx).
		Where("id = ?", batchID).
		First(&batch).Error; err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("❌ Error ambil analisis batch:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil analisis batch",
		})
	}
	fmt.Printf("🧠 Batch ID dari forecast: %s\n", batchID)
	fmt.Println("📝 Analisis batch ditemukan.")

	return c.Status(http.StatusOK).JSON(SalesChartResponse{
		Historical: historical,
		Forecast:   forecasts,
		Analysis:   batch.Analysis,
	})
}
