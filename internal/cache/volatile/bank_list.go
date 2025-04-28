package volatile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func refreshBankList() {
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
		log.Println("[CACHE] 🔄 Refresh BankList (Key ID, Value JSON {name, code})")

		var banks []models.BankList
		if err := database.DB.Order("id ASC").Find(&banks).Error; err != nil {
			log.Println("[CACHE] Gagal refresh BankList:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, bank := range banks {
			key := fmt.Sprintf("bank:%d", bank.ID)

			value := map[string]string{
				"name": bank.Name,
				"code": bank.Code,
			}

			jsonValue, err := json.Marshal(value)
			if err != nil {
				log.Println("[CACHE] Gagal encode JSON BankList:", err)
				continue
			}

			if err := redisClient.Set(ctx, key, jsonValue, 24*time.Hour).Err(); err != nil {
				log.Println("[CACHE] Gagal set BankList:", err)
			}
		}

		<-ticker.C
	}
}
