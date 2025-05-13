package volatile

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/config"
)

func GetVolatileRedisClient(ctx context.Context) (*redis.Client, error) {
	prefix, err := GetVolatileRedisPrefix()
	if err != nil {
		return nil, err
	}
	return config.GetRedisClient(prefix)
}
