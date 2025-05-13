package cloudinary

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func ResolveCloudinaryPrivatePrefix() (string, error) {
	ctx := context.Background()

	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err == nil {
		data, err := rdb.SMembers(ctx, "Private Image").Result()
		if err == nil && len(data) > 0 {
			for _, jsonStr := range data {
				var cs models.CloudService
				if err := json.Unmarshal([]byte(jsonStr), &cs); err != nil {
					log.Println("[CACHE] ⚠️ Gagal unmarshal CloudService:", err)
					continue
				}

				if cfg, err := GetCloudinaryConfig(cs.EnvKeyPrefix); err == nil && cfg.CloudName != "" {
					return cs.EnvKeyPrefix, nil
				}
			}
		}
	}

	var cs models.CloudService
	err = database.DB.
		Where("usage_for = ?", "Private Image").
		Order("storage_usage ASC").
		First(&cs).Error
	if err != nil {
		return "", errors.New("gagal ambil Cloudinary service dari DB")
	}

	cfg, err := GetCloudinaryConfig(cs.EnvKeyPrefix)
	if err != nil || cfg.CloudName == "" {
		return "", errors.New("cloudinary config tidak valid")
	}

	return cs.EnvKeyPrefix, nil
}
