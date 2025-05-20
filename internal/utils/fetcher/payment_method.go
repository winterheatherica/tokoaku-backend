package fetcher

import (
	"context"
	"log"
	"strconv"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func GetAllPaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal ambil Redis client:", err)
		rdb = nil
	}

	if rdb != nil {
		data, err := rdb.HGetAll(ctx, "payment_method:id").Result()
		if err == nil && len(data) > 0 {
			var methods []models.PaymentMethod
			for idStr, name := range data {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("[CACHE] ⚠️ Gagal parsing ID %s: %v", idStr, err)
					continue
				}
				methods = append(methods, models.PaymentMethod{
					ID:   uint(id),
					Name: name,
				})
			}
			log.Printf("[CACHE] ✅ Berhasil ambil %d payment method dari Redis", len(methods))
			return methods, nil
		}
		log.Println("[CACHE] ℹ️ Redis kosong atau error, fallback ke DB")
	}

	var methods []models.PaymentMethod
	if err := database.DB.Order("id ASC").Find(&methods).Error; err != nil {
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d payment method dari database", len(methods))

	if rdb != nil {
		entries := make(map[string]string)
		for _, m := range methods {
			entries[strconv.Itoa(int(m.ID))] = m.Name
		}
		if err := rdb.HSet(ctx, "payment_method:id", entries).Err(); err != nil {
			log.Println("[CACHE] ⚠️ Gagal menyimpan payment method ke Redis:", err)
		} else {
			log.Printf("[CACHE] ✅ Payment method disimpan ke Redis (%d item)", len(entries))
		}
	}

	return methods, nil
}
