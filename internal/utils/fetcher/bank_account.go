package fetcher

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

var (
	platformUserID string
	loadOnce       sync.Once
)

func loadPlatformUserID() {
	loadOnce.Do(func() {
		id := os.Getenv("PLATFORM_ID")
		if id == "" {
			log.Fatalf("❌ PLATFORM_ID not set in environment")
		}
		platformUserID = id
		log.Printf("✅ PLATFORM_ID loaded: %s", platformUserID)
	})
}

func GetAllPlatformBankAccounts(ctx context.Context) ([]models.BankAccount, error) {
	loadPlatformUserID()

	var accounts []models.BankAccount
	err := database.DB.WithContext(ctx).
		Preload("Bank").
		Where("user_id = ?", platformUserID).
		Find(&accounts).Error

	if err != nil {
		log.Printf("[DB] ❌ Gagal ambil bank account platform (user_id = %s): %v", platformUserID, err)
		return nil, err
	}

	log.Printf("[DB] ✅ Ditemukan %d bank account milik platform", len(accounts))
	return accounts, nil
}

func GetActiveBankAccountByUserID(ctx context.Context, userID string) (*models.BankAccount, error) {
	var account models.BankAccount

	err := database.DB.WithContext(ctx).
		Preload("Bank").
		Preload("User").
		Where("user_id = ? AND is_active = TRUE", userID).
		First(&account).Error

	if err != nil {
		log.Printf("[DB] ❌ Gagal ambil bank account aktif untuk user_id %s: %v", userID, err)
		return nil, err
	}

	log.Printf("[DB] ✅ Bank account aktif ditemukan untuk user_id %s: %s", userID, account.AccountNumber)

	return &account, nil
}
