package utils

import (
	"log"
	"sync"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

var (
	persistentRedisPrefix string
	persistentOnce        sync.Once
)

func LoadPersistentRedisPrefix() error {
	var err error

	persistentOnce.Do(func() {
		var cloudService models.CloudService
		if e := database.DB.
			Where("usage_for = ?", "Persistent Cache").
			Limit(1).
			First(&cloudService).Error; e != nil {
			log.Println("Gagal ambil persistent cache prefix:", e)
			err = e
			return
		}

		persistentRedisPrefix = cloudService.EnvKeyPrefix
		log.Println("[UTILS] Persistent Redis Prefix loaded:", persistentRedisPrefix)
	})

	return err
}

func GetPersistentRedisPrefix() (string, error) {
	if persistentRedisPrefix == "" {
		err := LoadPersistentRedisPrefix()
		if err != nil {
			return "", err
		}
	}
	return persistentRedisPrefix, nil
}
