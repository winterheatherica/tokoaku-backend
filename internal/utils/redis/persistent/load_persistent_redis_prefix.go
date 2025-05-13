package persistent

import (
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func LoadPersistentRedisPrefix() error {
	var err error

	persistentOnce.Do(func() {
		var cloudService models.CloudService
		if e := database.DB.
			Where("usage_for = ?", "Persistent Cache").
			Limit(1).
			First(&cloudService).Error; e != nil {
			log.Println("[CACHE] ❌ Gagal ambil Persistent Redis prefix:", e)
			err = e
			return
		}

		persistentRedisPrefix = cloudService.EnvKeyPrefix
	})

	return err
}
