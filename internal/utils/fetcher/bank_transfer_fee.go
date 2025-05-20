package fetcher

import (
	"context"
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func GetCheapestBankTransferFee(ctx context.Context, fromBankIDs []uint, toBankID uint) (*models.BankTransferFee, error) {
	var fees []models.BankTransferFee

	err := database.DB.WithContext(ctx).
		Preload("FromBank").
		Preload("ToBank").
		Preload("Fee").
		Where("from_bank_id IN ? AND to_bank_id = ?", fromBankIDs, toBankID).
		Find(&fees).Error

	if err != nil {
		log.Printf("[DB] ❌ Gagal ambil fee transfer dari bank_ids %+v ke %d: %v", fromBankIDs, toBankID, err)
		return nil, err
	}

	if len(fees) == 0 {
		log.Printf("[DB] ⚠️ Tidak ada fee transfer ditemukan dari bank_ids %+v ke %d", fromBankIDs, toBankID)
		return nil, nil
	}

	cheapest := &fees[0]
	for _, f := range fees[1:] {
		if f.Fee.Fee < cheapest.Fee.Fee {
			cheapest = &f
		}
	}

	log.Printf("[DB] ✅ Fee termurah: %d (%s)", cheapest.Fee.Fee, cheapest.Fee.ServiceName)
	return cheapest, nil
}
