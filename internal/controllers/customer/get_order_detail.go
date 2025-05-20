package customer

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetOrderDetail(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	userUID := c.Locals("uid").(string)

	orderUint, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Order ID tidak valid.")
	}

	var order models.Order
	if err := database.DB.
		Preload("Address").
		Preload("PaymentMethod").
		Where("id = ? AND customer_id = ?", uint(orderUint), userUID).
		First(&order).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Order tidak ditemukan.")
	}

	var orderLogs []models.OrderLog
	database.DB.
		Preload("Status").
		Where("order_id = ?", order.ID).
		Order("created_at ASC").
		Find(&orderLogs)

	var lastStatusID uint = 0
	if len(orderLogs) > 0 {
		lastStatusID = orderLogs[len(orderLogs)-1].StatusID
	}

	var orderPromo models.OrderPromo
	var promo *models.Promo
	if err := database.DB.
		Where("order_id = ?", order.ID).
		First(&orderPromo).Error; err == nil {
		var p models.Promo
		if err := database.DB.
			Where("id = ?", orderPromo.PromoID).
			First(&p).Error; err == nil {
			promo = &p
		}
	}

	var shippings []models.OrderShipping
	database.DB.
		Preload("Seller").
		Where("order_id = ?", order.ID).
		Find(&shippings)

	type ItemWithPrice struct {
		models.OrderItem
		Price *fetcher.PriceWithDiscountResponse `json:"price"`
	}

	type ShippingWithDetails struct {
		models.OrderShipping
		Seller models.User                  `json:"seller"`
		Status []models.OrderShippingStatus `json:"statuses"`
		Items  []ItemWithPrice              `json:"items"`
	}

	var shippingDetails []ShippingWithDetails

	for _, ship := range shippings {
		var statuses []models.OrderShippingStatus
		database.DB.
			Preload("Status").
			Where("order_shipping_id = ?", ship.ID).
			Order("created_at ASC").
			Find(&statuses)

		var items []models.OrderItem
		database.DB.
			Preload("ProductVariant").
			Preload("ProductVariant.Product").
			Where("order_shipping_id = ?", ship.ID).
			Find(&items)

		var enrichedItems []ItemWithPrice
		for _, item := range items {
			price, _ := fetcher.GetHistoricalPriceWithDiscount(c.Context(), item.ProductVariantID, order.CreatedAt)
			enrichedItems = append(enrichedItems, ItemWithPrice{
				OrderItem: item,
				Price:     price,
			})
		}

		shippingDetails = append(shippingDetails, ShippingWithDetails{
			OrderShipping: ship,
			Seller:        ship.Seller,
			Status:        statuses,
			Items:         enrichedItems,
		})
	}

	return c.JSON(fiber.Map{
		"order":           order,
		"order_logs":      orderLogs,
		"promo":           promo,
		"order_shippings": shippingDetails,
		"last_status":     lastStatusID,
	})
}
