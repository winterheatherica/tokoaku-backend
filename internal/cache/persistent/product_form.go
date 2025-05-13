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

func storeProductFormToRedis(ctx context.Context, rdb *redis.Client, form models.ProductForm) {
	idKey := "product_form:id"
	formKey := "product_form:form"

	if err := rdb.HSet(ctx, idKey, fmt.Sprintf("%d", form.ID), form.Form).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s → ID %d: %v", idKey, form.ID, err)
	}

	idJSON, err := json.Marshal(form.ID)
	if err != nil {
		log.Printf("[CACHE] ❌ Gagal Marshal ID JSON untuk %s: %v", form.Form, err)
		return
	}

	if err := rdb.HSet(ctx, formKey, form.Form, idJSON).Err(); err != nil {
		log.Printf("[CACHE] ❌ Gagal HSet %s → Form '%s': %v", formKey, form.Form, err)
	}
}

func refreshProductForms() {
	go func() {
		log.Println("[CACHE] ▶️  Memulai goroutine refreshProductForms (Persistent)")

		ctx := context.Background()

		rdb, err := persistent.GetPersistentRedisClient(ctx)
		if err != nil {
			log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
			return
		}

		ticker := time.NewTicker(cache.TickInterval24h)
		defer ticker.Stop()

		for {
			log.Println("[CACHE] 🔄 Refresh ProductForms (Persistent)")

			var productForms []models.ProductForm
			if err := database.DB.Order("id ASC").Find(&productForms).Error; err != nil {
				log.Println("[CACHE] ❌ Gagal mengambil data ProductForms:", err)
				time.Sleep(cache.SleepOnError)
				continue
			}

			for _, form := range productForms {
				storeProductFormToRedis(ctx, rdb, form)
			}

			log.Printf("[CACHE] ✅ Berhasil refresh %d product form ke Redis (Persistent)", len(productForms))
			<-ticker.C
		}
	}()
}
