package customer

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func ptrUint(u uint) *uint {
	return &u
}

func getVariantImageURL(variant models.ProductVariant) string {
	if len(variant.ProductVariantImages) > 0 {
		for _, img := range variant.ProductVariantImages {
			if img.IsVariantCover {
				return img.ImageURL
			}
		}
		return variant.ProductVariantImages[0].ImageURL
	}
	return variant.Product.ImageCoverURL
}

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
	if err := database.DB.
		Preload("Status").
		Where("order_id = ?", order.ID).
		Order("created_at ASC").
		Find(&orderLogs).Error; err != nil {
		log.Printf("Gagal ambil order logs: %v", err)
	}
	lastStatusID := uint(0)
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
			Preload("ValueType").
			Where("id = ?", orderPromo.PromoID).
			First(&p).Error; err == nil {
			promo = &p
		}
	}

	var shippings []models.OrderShipping
	if err := database.DB.
		Preload("Seller").
		Preload("ShippingOption").
		Preload("OrderItems.ProductVariant.Product").
		Preload("OrderItems.ProductVariant.ProductVariantImages").
		Where("order_id = ?", order.ID).
		Find(&shippings).Error; err != nil {
		log.Printf("Gagal preload shippings: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data pengiriman.")
	}

	platformAccounts, err := fetcher.GetAllPlatformBankAccounts(c.Context())
	if err != nil {
		log.Printf("Gagal ambil akun bank platform: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil akun bank platform")
	}
	fromBankIDs := []uint{}
	for _, acc := range platformAccounts {
		fromBankIDs = append(fromBankIDs, acc.BankID)
	}

	const platformFeePerSeller uint = 2000
	var shippingDetails []fiber.Map
	var totalSubtotal, totalTransferFee, totalShippingFee uint

	for _, ship := range shippings {
		var statuses []models.OrderShippingStatus
		if err := database.DB.
			Preload("Status").
			Where("order_shipping_id = ?", ship.ID).
			Order("created_at ASC").
			Find(&statuses).Error; err != nil {
			log.Printf("Gagal ambil status shipping %d: %v", ship.ID, err)
		}

		var itemDetails []fiber.Map
		sellerSubtotal := uint(0)

		for _, item := range ship.OrderItems {
			price, err := fetcher.GetHistoricalPriceWithDiscount(c.Context(), item.ProductVariantID, order.CreatedAt)
			if err != nil || price == nil {
				price = &fetcher.PriceWithDiscountResponse{
					ProductVariantID: item.ProductVariantID,
					OriginalPrice:    ptrUint(0),
					FinalPrice:       ptrUint(0),
					Discounts:        []fetcher.DiscountDTO{},
				}
			}

			originalPrice := uint(0)
			finalPrice := uint(0)
			if price.OriginalPrice != nil {
				originalPrice = *price.OriginalPrice
			}
			if price.FinalPrice != nil {
				finalPrice = *price.FinalPrice
			}

			subtotal := item.Quantity * finalPrice
			sellerSubtotal += subtotal

			imageURL := ""
			if item.ProductVariant != nil {
				imageURL = getVariantImageURL(*item.ProductVariant)
			}

			itemDetails = append(itemDetails, fiber.Map{
				"item": fiber.Map{
					"product_variant_id": item.ProductVariantID,
					"quantity":           item.Quantity,
				},
				"price": fiber.Map{
					"product_variant_id": price.ProductVariantID,
					"original_price":     originalPrice,
					"final_price":        finalPrice,
					"discounts":          price.Discounts,
				},
				"subtotal":  subtotal,
				"image_url": imageURL,
			})
		}

		transferFee := uint(0)
		sellerBank, err := fetcher.GetActiveBankAccountByUserID(c.Context(), ship.SellerID)
		if err == nil && sellerBank != nil {
			feeData, _ := fetcher.GetCheapestBankTransferFee(c.Context(), fromBankIDs, sellerBank.BankID)
			if feeData != nil {
				transferFee = feeData.Fee.Fee
			}
		}

		shippingFee := uint(0)
		if ship.ShippingOption.ID != 0 {
			shippingFee = ship.ShippingOption.Fee
		}

		totalSubtotal += sellerSubtotal
		totalTransferFee += transferFee
		totalShippingFee += shippingFee

		shippingDetails = append(shippingDetails, fiber.Map{
			"shipping":     ship,
			"seller":       ship.Seller,
			"statuses":     statuses,
			"items":        itemDetails,
			"transfer_fee": transferFee,
			"shipping_fee": shippingFee,
			"platform_fee": platformFeePerSeller,
		})
	}

	totalPlatformFee := uint(len(shippings)) * platformFeePerSeller

	var promoDiscount uint
	if promo != nil {
		totalBeforePromo := totalSubtotal + totalTransferFee + totalShippingFee + totalPlatformFee
		switch promo.ValueTypeID {
		case 1:
			promoDiscount = uint(float64(totalBeforePromo) * (float64(promo.Value) / 100))
		case 2:
			promoDiscount = promo.Value
		default:
			promoDiscount = 0
		}
		if promo.MaxValue > 0 && promoDiscount > promo.MaxValue {
			promoDiscount = promo.MaxValue
		}
	}

	grandTotal := totalSubtotal + totalTransferFee + totalShippingFee + totalPlatformFee
	if promoDiscount > 0 && grandTotal > promoDiscount {
		grandTotal -= promoDiscount
	}

	paymentMethods, _ := fetcher.GetAllPaymentMethods(c.Context())
	var paymentMethodResponse []fiber.Map
	for _, method := range paymentMethods {
		paymentMethodResponse = append(paymentMethodResponse, fiber.Map{
			"id":   method.ID,
			"name": method.Name,
		})
	}

	return c.JSON(fiber.Map{
		"order":           order,
		"order_logs":      orderLogs,
		"promo":           promo,
		"order_shippings": shippingDetails,
		"last_status":     lastStatusID,
		"checkout_summary": fiber.Map{
			"total": fiber.Map{
				"subtotal":       totalSubtotal,
				"transfer_fee":   totalTransferFee,
				"shipping_fee":   totalShippingFee,
				"platform_fee":   totalPlatformFee,
				"promo_discount": promoDiscount,
				"grand_total":    grandTotal,
			},
			"payment_methods": paymentMethodResponse,
		},
	})
}
