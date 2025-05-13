package fetcher

import (
	"context"
	"log"
	"strconv"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func GetAllProductForms(ctx context.Context) ([]models.ProductForm, error) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal ambil Redis client:", err)
		rdb = nil
	}

	if rdb != nil {
		data, err := rdb.HGetAll(ctx, "product_form:id").Result()
		if err == nil && len(data) > 0 {
			var forms []models.ProductForm
			for idStr, formStr := range data {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("[CACHE] ⚠️ Gagal parsing ID %s: %v", idStr, err)
					continue
				}
				forms = append(forms, models.ProductForm{
					ID:   uint(id),
					Form: formStr,
				})
			}
			log.Printf("[CACHE] ✅ Berhasil ambil %d product form dari Redis", len(forms))
			return forms, nil
		}
		log.Println("[CACHE] ℹ️ Redis kosong atau error, fallback ke DB")
	}

	var forms []models.ProductForm
	if err := database.DB.Order("form ASC").Find(&forms).Error; err != nil {
		return nil, err
	}
	log.Printf("[DB] ✅ Berhasil ambil %d product form dari database", len(forms))

	if rdb != nil {
		entries := make(map[string]string)
		for _, form := range forms {
			entries[strconv.Itoa(int(form.ID))] = form.Form
		}
		if err := rdb.HSet(ctx, "product_form:id", entries).Err(); err != nil {
			log.Println("[CACHE] ⚠️ Gagal menyimpan product form ke Redis:", err)
		} else {
			log.Printf("[CACHE] ✅ Product form disimpan ke Redis (%d item)", len(entries))
		}
	}

	return forms, nil
}
