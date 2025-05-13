package persistent

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func storeBankTransferFeeToRedis(ctx context.Context, rdb *redis.Client, fee models.BankTransferFee) {
	field := fmt.Sprintf("%d:%d", fee.FromBankID, fee.ToBankID)
	if err := rdb.HSet(ctx, "bank_transfer_fee", field, fee.Fee.Fee).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet BankTransferFee %s: %v", field, err)
	}
}

func refreshBankTransferFees() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine RefreshBankTransferFees (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval24h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh BankTransferFees (Persistent)")

			var fees []models.BankTransferFee
			err := database.DB.Preload("Fee").
				Order("from_bank_id ASC, to_bank_id ASC").
				Find(&fees).Error
			if err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data BankTransferFees:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, fee := range fees {
				storeBankTransferFeeToRedis(ctx, rdb, fee)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d data fee transfer bank ke Redis (Persistent)", len(fees))
			<-ticker.C
		}
	}()
}
