package admin

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func UploadSalesCSV(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File tidak ditemukan",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuka file",
		})
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil || len(records) <= 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File CSV tidak valid atau kosong",
		})
	}

	skipped := 0
	inserted := 0

	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}

		dateStr := strings.TrimSpace(row[0])
		salesStr := strings.TrimSpace(row[len(row)-1])

		var date time.Time
		if strings.Contains(dateStr, "-") {
			date, err = time.Parse("2006-01-02", dateStr)
		} else {
			date, err = time.Parse("1/2/2006", dateStr)
		}
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Format tanggal salah di baris %d: %v", i+1, dateStr),
			})
		}

		totalSales, err := strconv.ParseInt(salesStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Total sales tidak valid di baris %d: %v", i+1, salesStr),
			})
		}

		var existing models.SalesData
		err = database.DB.First(&existing, "date = ?", date).Error
		if err == nil {
			skipped++
			continue
		}

		data := models.SalesData{
			Date:       date,
			TotalSales: totalSales,
		}
		if err := database.DB.Create(&data).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Gagal menyimpan data di baris %d", i+1),
			})
		}
		inserted++
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":  "Upload selesai",
		"inserted": inserted,
		"skipped":  skipped,
	})
}
