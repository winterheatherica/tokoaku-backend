package redis

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/config"
)

var Client *redis.Client
var Ctx = context.Background()

func Connect() {
	opts := &redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
	}

	if config.Redis.UseTLS {
		opts.TLSConfig = &tls.Config{}
	}

	Client = redis.NewClient(opts)

	if _, err := Client.Ping(Ctx).Result(); err != nil {
		log.Fatal("Redis gagal connect:", err)
	}

	log.Printf("[SERVICE]: ⚙️  Redis connected to %s (TLS: %v)", config.Redis.Host, config.Redis.UseTLS)
}
