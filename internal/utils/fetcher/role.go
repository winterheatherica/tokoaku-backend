package fetcher

import (
	"context"
	"log"
	"strconv"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func GetRoleNameByID(roleID int) (string, error) {
	ctx := context.Background()
	cacheKey := "role:" + strconv.Itoa(roleID)

	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		log.Println("[CACHE] ❌ Gagal mendapatkan Redis client:", err)
	} else {
		roleName, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			log.Printf("[CACHE] ✅ Role ID %d ditemukan: %s", roleID, roleName)
			return roleName, nil
		}
		log.Printf("[CACHE] ℹ️ Role ID %d tidak ditemukan di Redis, fallback ke DB", roleID)
	}

	var role models.Role
	if err := database.DB.First(&role, "id = ?", roleID).Error; err != nil {
		log.Printf("[DB] ❌ Role ID %d tidak ditemukan: %v", roleID, err)
		return "", err
	}

	if rdb != nil {
		if err := rdb.Set(ctx, cacheKey, role.Name, 0).Err(); err != nil {
			log.Printf("[CACHE] ⚠️ Gagal menyimpan role ID %d ke Redis: %v", roleID, err)
		} else {
			log.Printf("[CACHE] ✅ Role ID %d disimpan ke Redis", roleID)
		}
	}

	return role.Name, nil
}
