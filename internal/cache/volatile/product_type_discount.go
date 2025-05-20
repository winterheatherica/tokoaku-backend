package volatile

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/cache"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func storeDiscountProductTypeToRedis(ctx context.Context, rdb *redis.Client, discountMap map[uint][]uint) {
	exp := cache.TickInterval30m

	for discountID, productTypeIDs := range discountMap {
		if len(productTypeIDs) == 0 {
			continue
		}

		key := fmt.Sprintf("product_type_discount:%d", discountID)

		members := make([]interface{}, len(productTypeIDs))
		for i, id := range productTypeIDs {
			members[i] = id
		}

		if err := rdb.SAdd(ctx, key, members...).Err(); err != nil {
			log.Printf("[CACHE] ❌ Gagal SADD %s: %v", key, err)
			continue
		}

		if err := rdb.Expire(ctx, key, exp).Err(); err != nil {
			log.Printf("[CACHE] ❌ Gagal set TTL untuk %s: %v", key, err)
		}
	}
}

func refreshProductTypeDiscounts() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshProductTypeDiscounts (Volatile)")

		ctx := context.Background()
		rdb, err := volatile.GetVolatileRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval30m)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh ProductTypeDiscount (Volatile)")

			var productTypeDiscounts []models.ProductTypeDiscount
			if err := database.DB.
				Preload("Discount").
				Order("discount_id ASC, product_type_id ASC").
				Find(&productTypeDiscounts).Error; err != nil {

				log.Println("[CACHE] ❌ Gagal mengambil data ProductTypeDiscounts:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			discountMap := make(map[uint][]uint)
			now := time.Now()

			for _, ptd := range productTypeDiscounts {
				if ptd.Discount.StartAt.Before(now) && ptd.Discount.EndAt.After(now) {
					discountMap[ptd.DiscountID] = append(discountMap[ptd.DiscountID], ptd.ProductTypeID)
				}
			}

			storeDiscountProductTypeToRedis(ctx, rdb, discountMap)

			log.Printf("[CACHE] ✅ Berhasil refresh %d relasi diskon → tipe produk ke Redis (Volatile)", len(productTypeDiscounts))
			<-ticker.C
		}
	}()
}
