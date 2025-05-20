package customer

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func PurchaseOrderDemo(c *fiber.Ctx) error {
	orderIDStr := c.Params("order_id")
	userUID := c.Locals("uid").(string)

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Order ID tidak valid.")
	}

	var order models.Order
	if err := database.DB.
		Where("id = ? AND customer_id = ?", uint(orderID), userUID).
		First(&order).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Order tidak ditemukan.")
	}

	var lastLog models.OrderLog
	if err := database.DB.
		Where("order_id = ?", order.ID).
		Order("created_at DESC").
		First(&lastLog).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Status order tidak ditemukan.")
	}

	if lastLog.StatusID != 11 {
		return fiber.NewError(fiber.StatusBadRequest, "Order tidak dapat dibayar karena status terakhir bukan pending.")
	}

	if order.PaymentMethodID != 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Metode pembayaran tidak valid untuk pembayaran langsung.")
	}

	log := models.OrderLog{
		OrderID:  order.ID,
		StatusID: 12,
	}
	if err := database.DB.Create(&log).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mencatat status pembayaran.")
	}

	return c.JSON(fiber.Map{
		"message":   "Pembayaran berhasil dicatat.",
		"order_id":  order.ID,
		"status_id": log.StatusID,
	})
}
