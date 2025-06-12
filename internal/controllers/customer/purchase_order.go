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

	var shippings []models.OrderShipping
	if err := database.DB.Where("order_id = ?", order.ID).Find(&shippings).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data shipping untuk order.")
	}

	var shippingIDs []uint
	for _, s := range shippings {
		shippingIDs = append(shippingIDs, s.ID)
	}

	var items []models.OrderItem
	if err := database.DB.Where("order_shipping_id IN ?", shippingIDs).Find(&items).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil item dalam order.")
	}

	var variantIDs []string
	for _, item := range items {
		variantIDs = append(variantIDs, item.ProductVariantID)
	}

	if err := database.DB.
		Model(&models.Cart{}).
		Where("customer_id = ? AND product_variant_id IN ? AND is_converted = false", userUID, variantIDs).
		Update("is_converted", true).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengubah status keranjang.")
	}

	return c.JSON(fiber.Map{
		"message":   "Pembayaran berhasil dicatat.",
		"order_id":  order.ID,
		"status_id": log.StatusID,
	})
}
