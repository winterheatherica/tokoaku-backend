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

func GetRoleNameByID(roleID int) (string, error) {
	ctx := context.Background()
	key := "role:" + strconv.Itoa(roleID)

	prefix, err := GetVolatileRedisPrefix()
	if err != nil {
		return "", err
	}

	redisClient, err := redis.GetRedisClient(prefix)
	if err != nil {
		log.Println("Gagal ambil Redis client:", err)
		return "", err
	}

	roleName, err := redisClient.Get(ctx, key).Result()
	if err == nil {
		return roleName, nil
	}

	var role models.Role
	if err := database.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return "", err
	}

	if err := redisClient.Set(ctx, key, role.Name, 24*time.Hour).Err(); err != nil {
		log.Println("⚠️ Gagal cache role ke redis:", err)
	}

	return role.Name, nil
}
