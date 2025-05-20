package customer

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

type MinimalShipping struct {
	SellerID         string `json:"seller_id"`
	ShippingOptionID uint   `json:"shipping_option_id"`
	BankAccountID    string `json:"bank_account_id"`
}

type OrderRequest struct {
	PaymentMethodID uint              `json:"payment_method_id"`
	AddressID       string            `json:"address_id"`
	PromoID         *uint             `json:"promo_id"`
	OrderShippings  []MinimalShipping `json:"order_shippings"`
}

func CreateOrder(c *fiber.Ctx) error {
	rawUID := c.Locals("uid")
	userUID, ok := rawUID.(string)
	if !ok || userUID == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "User belum login.")
	}

	user, err := fetcher.GetUserByUID(c.Context(), userUID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User tidak ditemukan.")
	}

	var req OrderRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Request tidak valid: "+err.Error())
	}

	addressUUID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Address ID tidak valid.")
	}

	var address models.Address
	if err := database.DB.
		Where("id = ? AND user_id = ?", addressUUID, userUID).
		First(&address).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Alamat tidak ditemukan.")
	}

	var selectedCarts []models.Cart
	if err := database.DB.
		Preload("ProductVariant.Product").
		Preload("ProductVariant").
		Where("customer_id = ? AND is_selected = true AND is_converted = false", userUID).
		Find(&selectedCarts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil cart.")
	}

	productForms, err := fetcher.GetAllProductForms(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil form produk")
	}
	productFormMap := make(map[uint]string)
	for _, pf := range productForms {
		productFormMap[pf.ID] = pf.Form
	}

	for i := range req.OrderShippings {
		ship := &req.OrderShippings[i]
		hasPhysical := false
		for _, item := range selectedCarts {
			if item.ProductVariant.Product.SellerID == ship.SellerID {
				form := productFormMap[item.ProductVariant.Product.ProductFormID]
				if form == "Physical" {
					hasPhysical = true
					break
				}
			}
		}
		if !hasPhysical {
			ship.ShippingOptionID = 100
		}
	}

	var totalSubtotal uint = 0
	for _, item := range selectedCarts {
		priceInfo, _ := fetcher.GetPriceWithDiscountForUI(c.Context(), item.ProductVariantID)
		if priceInfo != nil && priceInfo.FinalPrice != nil {
			totalSubtotal += item.Quantity * *priceInfo.FinalPrice
		}
	}

	platformAccounts, _ := fetcher.GetAllPlatformBankAccounts(c.Context())
	fromBankIDs := []uint{}
	for _, acc := range platformAccounts {
		fromBankIDs = append(fromBankIDs, acc.BankID)
	}

	var totalTransferFee uint = 0
	feePerSeller := make(map[string]*models.BankTransferFee)
	for _, ship := range req.OrderShippings {
		if _, exists := feePerSeller[ship.SellerID]; !exists {
			sellerBank, err := fetcher.GetActiveBankAccountByUserID(c.Context(), ship.SellerID)
			if err == nil && sellerBank != nil {
				fee, _ := fetcher.GetCheapestBankTransferFee(c.Context(), fromBankIDs, sellerBank.BankID)
				feePerSeller[ship.SellerID] = fee
			}
		}
		if fee := feePerSeller[ship.SellerID]; fee != nil {
			totalTransferFee += fee.Fee.Fee
		}
	}

	var totalShippingFee uint = 0
	for _, ship := range req.OrderShippings {
		if ship.ShippingOptionID == 100 {
			continue
		}
		opts, _ := fetcher.GetSellerShippingOptions(c.Context(), ship.SellerID)
		for _, opt := range opts {
			if opt.ShippingOptionID == ship.ShippingOptionID {
				totalShippingFee += uint(opt.Fee)
				break
			}
		}
	}

	var promoDiscount uint = 0
	if req.PromoID != nil {
		userPromo, err := fetcher.GetUserPromoByPromoID(c.Context(), user.ID, *req.PromoID)
		if err == nil && userPromo != nil {
			discount, err := fetcher.CalculatePromoDiscount(userPromo.Promo, totalSubtotal)
			if err == nil {
				promoDiscount = discount
			}
		}
	}

	const platformFee uint = 2000
	grandTotal := totalSubtotal + totalTransferFee + totalShippingFee + platformFee - promoDiscount

	order := models.Order{
		CustomerID:      userUID,
		PaymentMethodID: req.PaymentMethodID,
		AddressID:       addressUUID,
		TotalPrice:      grandTotal,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal simpan order: "+err.Error())
	}

	orderLog := models.OrderLog{
		OrderID:  order.ID,
		StatusID: 11,
	}
	database.DB.Create(&orderLog)

	if req.PromoID != nil {
		orderPromo := models.OrderPromo{
			OrderID: order.ID,
			PromoID: *req.PromoID,
		}
		database.DB.Create(&orderPromo)

		database.DB.Model(&models.UserPromo{}).
			Where("promo_id = ? AND customer_id = ?", *req.PromoID, userUID).
			Update("redeemed", true)
	}

	groupedCarts := make(map[string][]models.Cart)
	for _, item := range selectedCarts {
		sellerID := item.ProductVariant.Product.SellerID
		groupedCarts[sellerID] = append(groupedCarts[sellerID], item)
	}

	for _, ship := range req.OrderShippings {
		bankUUID, err := uuid.Parse(ship.BankAccountID)
		if err != nil {
			log.Printf("❌ BankAccountID invalid: %v", err)
			continue
		}

		var trackingPtr *string
		if ship.ShippingOptionID != 100 {
			t := "TRK-" + uuid.New().String()[:8]
			trackingPtr = &t
		}

		shipping := models.OrderShipping{
			OrderID:          order.ID,
			SellerID:         ship.SellerID,
			ShippingOptionID: ship.ShippingOptionID,
			BankAccountID:    bankUUID,
			TrackingNumber:   trackingPtr,
		}
		if err := database.DB.Create(&shipping).Error; err != nil {
			log.Printf("❌ Gagal simpan OrderShipping: %v", err)
			continue
		}

		database.DB.Create(&models.OrderShippingStatus{
			OrderShippingID: shipping.ID,
			StatusID:        21,
		})

		items := groupedCarts[ship.SellerID]
		for _, item := range items {
			orderItem := models.OrderItem{
				OrderShippingID:  shipping.ID,
				ProductVariantID: item.ProductVariantID,
				Quantity:         item.Quantity,
			}

			if err := database.DB.Create(&orderItem).Error; err != nil {
				log.Printf("❌ Gagal simpan OrderItem: %v", err)
				continue
			}
		}
	}

	if err := database.DB.
		Model(&models.Cart{}).
		Where("customer_id = ? AND is_selected = true AND is_converted = false", userUID).
		Update("is_converted", true).Error; err != nil {
		log.Printf("❌ Gagal update cart jadi converted: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Order berhasil dibuat",
		"order_id":       order.ID,
		"total_price":    grandTotal,
		"subtotal":       totalSubtotal,
		"transfer_fee":   totalTransferFee,
		"shipping_fee":   totalShippingFee,
		"promo_discount": promoDiscount,
		"platform_fee":   platformFee,
	})
}
