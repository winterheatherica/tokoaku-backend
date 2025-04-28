package utils

import (
	"log"
	"sync"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

var (
	volatileRedisPrefix string
	volatileOnce        sync.Once
)

func LoadVolatileRedisPrefix() error {
	var err error

	volatileOnce.Do(func() {
		var cloudService models.CloudService
		if e := database.DB.
			Where("usage_for = ?", "Volatile Cache").
			Limit(1).
			First(&cloudService).Error; e != nil {
			log.Println("Gagal ambil volatile cache prefix:", e)
			err = e
			return
		}

		volatileRedisPrefix = cloudService.EnvKeyPrefix
		log.Println("[UTILS] Volatile Redis Prefix loaded:", volatileRedisPrefix)
	})

	return err
}

func GetVolatileRedisPrefix() (string, error) {
	if volatileRedisPrefix == "" {
		err := LoadVolatileRedisPrefix()
		if err != nil {
			return "", err
		}
	}
	return volatileRedisPrefix, nil
}
