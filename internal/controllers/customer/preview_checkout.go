package customer

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func PreviewCheckout(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)

	var selectedCarts []models.Cart
	if err := database.DB.
		Preload("ProductVariant.Product.Seller").
		Preload("ProductVariant.Product").
		Preload("ProductVariant").
		Where("customer_id = ? AND is_selected = true AND is_converted = false", userUID).
		Find(&selectedCarts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil cart untuk checkout")
	}

	platformAccounts, err := fetcher.GetAllPlatformBankAccounts(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil akun bank platform")
	}
	fromBankIDs := []uint{}
	for _, acc := range platformAccounts {
		fromBankIDs = append(fromBankIDs, acc.BankID)
	}

	productForms, err := fetcher.GetAllProductForms(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data product form")
	}
	productFormMap := make(map[uint]string)
	for _, pf := range productForms {
		productFormMap[pf.ID] = pf.Form
	}

	promos, err := fetcher.GetAllActivePromos(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data promo aktif")
	}

	paymentMethods, err := fetcher.GetAllPaymentMethods(context.Background())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data metode pembayaran")
	}

	sellerGrouped := make(map[string][]fiber.Map)
	feePerSeller := make(map[string]*models.BankTransferFee)
	pricePerVariant := make(map[string]*fetcher.PriceWithDiscountResponse)
	sellerIDMap := make(map[string]string)

	bankAccountPerSeller := make(map[string]*models.BankAccount)

	for _, item := range selectedCarts {
		variantImages, err := fetcher.GetAllVariantImages(context.Background(), item.ProductVariantID)
		if err != nil {
			variantImages = []models.ProductVariantImage{}
		}

		imageURL := item.ProductVariant.Product.ImageCoverURL
		for _, img := range variantImages {
			if img.IsVariantCover {
				imageURL = img.ImageURL
				break
			}
		}

		images := []string{}
		for _, img := range variantImages {
			images = append(images, img.ImageURL)
		}

		priceInfo, _ := fetcher.GetPriceWithDiscountForUI(context.Background(), item.ProductVariantID)
		pricePerVariant[item.ProductVariantID] = priceInfo

		sellerName := "-"
		if item.ProductVariant.Product.Seller.Name != nil {
			sellerName = *item.ProductVariant.Product.Seller.Name
		}
		sellerID := item.ProductVariant.Product.SellerID
		sellerIDMap[sellerName] = sellerID

		if _, exists := feePerSeller[sellerName]; !exists {
			sellerBank, err := fetcher.GetActiveBankAccountByUserID(context.Background(), sellerID)
			if err == nil && sellerBank != nil {
				fee, _ := fetcher.GetCheapestBankTransferFee(context.Background(), fromBankIDs, sellerBank.BankID)
				feePerSeller[sellerName] = fee
				bankAccountPerSeller[sellerName] = sellerBank
			}
		}

		discounts := []fiber.Map{}
		if priceInfo != nil {
			for _, d := range priceInfo.Discounts {
				discounts = append(discounts, fiber.Map{
					"id":         d.ID,
					"name":       d.Name,
					"value":      d.Value,
					"value_type": d.ValueType,
					"sponsor":    d.Sponsor,
				})
			}
		}

		subtotal := uint(0)
		if priceInfo != nil && priceInfo.FinalPrice != nil {
			subtotal = item.Quantity * *priceInfo.FinalPrice
		}

		cartData := fiber.Map{
			"product_variant_id": item.ProductVariantID,
			"product_name":       item.ProductVariant.Product.Name,
			"variant_name":       item.ProductVariant.VariantName,
			"quantity":           item.Quantity,
			"image_url":          imageURL,
			"variant_images":     images,
			"added_at":           item.CreatedAt,
			"original_price":     priceInfo.OriginalPrice,
			"final_price":        priceInfo.FinalPrice,
			"discounts":          discounts,
			"subtotal":           subtotal,
			"product_form_id":    item.ProductVariant.Product.ProductFormID,
		}

		sellerGrouped[sellerName] = append(sellerGrouped[sellerName], cartData)
	}

	var totalSubtotal uint = 0
	var totalTransferFee uint = 0
	var response []fiber.Map

	for sellerName, items := range sellerGrouped {
		fee := feePerSeller[sellerName]
		var feeAmount uint = 0
		if fee != nil {
			feeAmount = fee.Fee.Fee
		}

		sellerID := sellerIDMap[sellerName]

		hasPhysical := false
		for _, item := range items {
			formID := item["product_form_id"].(uint)
			if formName, ok := productFormMap[formID]; ok && formName == "Physical" {
				hasPhysical = true
				break
			}
		}

		opts := []fiber.Map{}
		if hasPhysical {
			shippingOpts, _ := fetcher.GetSellerShippingOptions(context.Background(), sellerID)
			for _, o := range shippingOpts {
				opts = append(opts, fiber.Map{
					"shipping_option_id":   o.ShippingOptionID,
					"courier_name":         o.CourierName,
					"courier_service_name": o.CourierServiceName,
					"fee":                  o.Fee,
					"estimated_time":       o.EstimatedTime,
					"service_type":         o.ServiceType,
				})
			}
		}

		sellerSubtotal := uint(0)
		for _, item := range items {
			if sub, ok := item["subtotal"].(uint); ok {
				sellerSubtotal += sub
			}
		}

		totalSubtotal += sellerSubtotal
		totalTransferFee += feeAmount

		var bankAccountID string
		if bank := bankAccountPerSeller[sellerName]; bank != nil {
			bankAccountID = bank.ID.String()
		}

		response = append(response, fiber.Map{
			"seller_name":      sellerName,
			"seller_id":        sellerID,
			"cart_items":       items,
			"transfer_fee":     feeAmount,
			"shipping_options": opts,
			"bank_account_id":  bankAccountID,
		})
	}

	promoResponse := []fiber.Map{}
	for _, promo := range promos {
		promoResponse = append(promoResponse, fiber.Map{
			"id":              promo.ID,
			"name":            promo.Name,
			"code":            promo.Code,
			"description":     promo.Description,
			"value_type_id":   promo.ValueTypeID,
			"value_type":      promo.ValueType,
			"value":           promo.Value,
			"min_price_value": promo.MinPriceValue,
			"max_value":       promo.MaxValue,
			"start_at":        promo.StartAt,
			"end_at":          promo.EndAt,
		})
	}

	paymentMethodResponse := []fiber.Map{}
	for _, method := range paymentMethods {
		paymentMethodResponse = append(paymentMethodResponse, fiber.Map{
			"id":   method.ID,
			"name": method.Name,
		})
	}

	return c.JSON(fiber.Map{
		"items":           response,
		"promos":          promoResponse,
		"payment_methods": paymentMethodResponse,
		"total": fiber.Map{
			"subtotal":     totalSubtotal,
			"transfer_fee": totalTransferFee,
			"grand_total":  totalSubtotal + totalTransferFee,
		},
		"message": "Preview checkout berhasil",
	})
}
