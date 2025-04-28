package volatile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func refreshBankTransferFees() {
	ctx := context.Background()

	prefix, err := utils.GetVolatileRedisPrefix()
	if err != nil {
		log.Println("[CACHE] Gagal ambil volatile prefix:", err)
		return
	}

	redisClient, err := redis.GetRedisClient(prefix)
	if err != nil {
		log.Println("[CACHE] Gagal ambil Redis client:", err)
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		log.Println("[CACHE] 🔄 Refresh BankTransferFees")

		var fees []models.BankTransferFee
		if err := database.DB.Preload("Fee").Order("from_bank_id ASC, to_bank_id ASC").Find(&fees).Error; err != nil {
			log.Println("[CACHE] Gagal refresh BankTransferFees:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, fee := range fees {
			key := fmt.Sprintf("bank_transfer_fee:%d:%d", fee.FromBankID, fee.ToBankID)
			value := fee.Fee.Fee

			if err := redisClient.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
				log.Println("[CACHE] Gagal set BankTransferFee:", err)
			}
		}

		<-ticker.C
	}
}
