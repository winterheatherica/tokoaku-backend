package customer

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetSellerFee(c *fiber.Ctx) error {
	ctx := context.Background()
	customerID := c.Locals("uid").(string)

	var carts []models.Cart
	if err := database.DB.
		Preload("ProductVariant.Product").
		Where("customer_id = ?", customerID).
		Find(&carts).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal ambil keranjang")
	}

	grouped := make(map[string][]models.Cart)
	for _, cart := range carts {
		sellerID := cart.ProductVariant.Product.SellerID
		grouped[sellerID] = append(grouped[sellerID], cart)
	}

	platformBankAccounts, err := fetcher.GetAllPlatformBankAccounts(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal ambil akun bank platform")
	}

	var fromBankIDs []uint
	for _, acc := range platformBankAccounts {
		fromBankIDs = append(fromBankIDs, acc.BankID)
	}

	result := make([]fiber.Map, 0)

	for sellerID, sellerCart := range grouped {
		sellerBank, err := fetcher.GetActiveBankAccountByUserID(ctx, sellerID)
		if err != nil || sellerBank == nil {
			log.Printf("[SKIP] Tidak ditemukan akun bank aktif untuk seller %s", sellerID)
			continue
		}

		cheapestFee, err := fetcher.GetCheapestBankTransferFee(ctx, fromBankIDs, sellerBank.BankID)
		if err != nil {
			log.Printf("[SKIP] Tidak bisa ambil fee termurah ke seller %s: %v", sellerID, err)
			continue
		}

		result = append(result, fiber.Map{
			"seller_id":  sellerID,
			"fee":        cheapestFee.Fee.Fee,
			"service":    cheapestFee.Fee.ServiceName,
			"bank_name":  sellerBank.Bank.Name,
			"account_no": sellerBank.AccountNumber,
			"cart_items": sellerCart,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"transfer_fees": result,
	})
}
