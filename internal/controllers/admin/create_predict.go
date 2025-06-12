package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type flaskPredictionResponse struct {
	Predictions map[string]struct {
		TotalSales string `json:"total_sales"`
	} `json:"predictions"`
	Analysis string `json:"analysis"`
}

func PredictSales(c *fiber.Ctx) error {
	const WINDOW = 60
	const FORECAST_STEP = 30

	var salesData []models.SalesData
	if err := database.DB.
		Order("date DESC").
		Limit(WINDOW).
		Find(&salesData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data penjualan",
		})
	}

	if len(salesData) < WINDOW {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Data historis kurang dari %d hari", WINDOW),
		})
	}

	for i, j := 0, len(salesData)-1; i < j; i, j = i+1, j-1 {
		salesData[i], salesData[j] = salesData[j], salesData[i]
	}

	valid := false
	for _, s := range salesData {
		if s.TotalSales > 0 {
			valid = true
			break
		}
	}
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Semua data penjualan bernilai 0, tidak bisa memproses prediksi",
		})
	}

	type salesInput struct {
		Date       string `json:"date"`
		TotalSales int64  `json:"total_sales"`
	}
	var inputData []salesInput
	for _, s := range salesData {
		inputData = append(inputData, salesInput{
			Date:       s.Date.Format("2006-01-02"),
			TotalSales: s.TotalSales,
		})
	}

	requestBody := map[string]interface{}{"sales": inputData}
	jsonBody, _ := json.Marshal(requestBody)

	resp, err := http.Post(
		os.Getenv("MACHINE_LEARNING_BASE_URL")+"/predict-sales",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menghubungi service prediksi: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	var result flaskPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal membaca respons dari model prediksi",
		})
	}

	lastHistDate := salesData[len(salesData)-1].Date
	startForecast := lastHistDate.AddDate(0, 0, 1)
	endForecast := startForecast.AddDate(0, 0, FORECAST_STEP-1)

	batch := models.SalesForecastBatch{
		ID:        uuid.NewString(),
		StartDate: startForecast,
		EndDate:   endForecast,
		Analysis:  result.Analysis,
	}
	if err := database.DB.Create(&batch).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menyimpan batch prediksi",
		})
	}

	for dateStr, pred := range result.Predictions {
		dateParsed, _ := time.Parse("2006-01-02", dateStr)

		sanitized := strings.ReplaceAll(pred.TotalSales, ".", "")
		sanitized = strings.ReplaceAll(sanitized, "Rp", "")
		sanitized = strings.ReplaceAll(sanitized, "Rp.", "")
		sanitized = strings.TrimSpace(sanitized)

		var numeric int64
		fmt.Sscanf(sanitized, "%d", &numeric)

		forecast := models.SalesForecast{
			Date:           dateParsed,
			PredictedSales: numeric,
			BatchID:        batch.ID,
		}
		database.DB.Create(&forecast)
	}

	return c.JSON(fiber.Map{
		"message":    "Prediksi penjualan berhasil disimpan",
		"batch_id":   batch.ID,
		"start_date": batch.StartDate.Format("2006-01-02"),
		"end_date":   batch.EndDate.Format("2006-01-02"),
		"analysis":   batch.Analysis,
	})
}
