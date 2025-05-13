package persistent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func storeBankToRedis(ctx context.Context, rdb *redis.Client, bank models.BankList) {
	idKey := "bank:id"
	nameKey := "bank:name"

	value := map[string]string{
		"name": bank.Name,
		"code": bank.Code,
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal encode JSON bank ID %d: %v", bank.ID, err)
		return
	}

	if err := rdb.HSet(ctx, idKey, fmt.Sprintf("%d", bank.ID), jsonValue).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal menyimpan ke Redis %s: %v", idKey, err)
	}

	if err := rdb.HSet(ctx, nameKey, bank.Name, bank.ID).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal menyimpan ke Redis %s: %v", nameKey, err)
	}
}

func refreshBankList() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine RefreshBankList (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval24h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Memulai refresh data BankList (Persistent)")

			var banks []models.BankList
			if err := database.DB.Order("id ASC").Find(&banks).Error; err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data dari database:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, bank := range banks {
				storeBankToRedis(ctx, rdb, bank)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d data bank ke Redis (Persistent)", len(banks))
			<-ticker.C
		}
	}()
}
