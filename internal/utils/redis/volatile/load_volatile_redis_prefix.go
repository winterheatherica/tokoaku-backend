package volatile

import (
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func LoadVolatileRedisPrefix() error {
	var err error

	volatileOnce.Do(func() {
		var cloudService models.CloudService
		if e := database.DB.
			Where("usage_for = ?", "Volatile Cache").
			Limit(1).
			First(&cloudService).Error; e != nil {
			log.Println("[CACHE] ❌ Gagal ambil Volatile Redis prefix:", e)
			err = e
			return
		}

		volatileRedisPrefix = cloudService.EnvKeyPrefix
	})

	return err
}
