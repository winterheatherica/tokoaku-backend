package persistent

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/config"
)

func GetPersistentRedisClient(ctx context.Context) (*redis.Client, error) {
	prefix, err := GetPersistentRedisPrefix()
	if err != nil {
		return nil, err
	}
	return config.GetRedisClient(prefix)
}
