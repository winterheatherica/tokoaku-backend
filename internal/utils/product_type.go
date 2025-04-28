package utils

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
)

func GetAllProductTypesFromCacheOrDB() ([]models.ProductType, error) {
	ctx := context.Background()

	prefix, err := GetVolatileRedisPrefix()
	if err != nil {
		log.Println("[ERROR] Gagal ambil volatile prefix:", err)
		return nil, err
	}

	redisClient, err := redis.GetRedisClient(prefix)
	if err != nil {
		log.Println("[ERROR] Gagal ambil Redis client:", err)
		return nil, err
	}

	var productTypes []models.ProductType
	err = database.DB.Find(&productTypes).Error
	if err != nil {
		log.Println("[ERROR] Gagal ambil ProductTypes dari DB:", err)
		return nil, err
	}

	for _, productType := range productTypes {
		key := "product_type_id:" + strconv.Itoa(int(productType.ID))
		value := productType.Name

		err := redisClient.Set(ctx, key, value, 24*time.Hour).Err()
		if err != nil {
			log.Printf("[ERROR] Gagal cache product type ke Redis: %s -> %s | %v", key, value, err)
			continue
		}
		log.Printf("[CACHE] Cached product type: %s -> %s", key, value)
	}

	return productTypes, nil
}
