package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

type RedisDynamic struct {
	Client  *redis.Client
	UseTLS  bool
	Address string
}

var redisClients map[string]*RedisDynamic

func LoadAllRedisConfigs(db *gorm.DB) error {
	var services []models.CloudService

	err := db.Where("provider_id = ?", 12).Find(&services).Error
	if err != nil {
		return fmt.Errorf("failed to load Redis cloud services: %w", err)
	}

	redisClients = make(map[string]*RedisDynamic)

	for _, service := range services {
		rawURL := os.Getenv(fmt.Sprintf("%s_REDIS_URL", service.EnvKeyPrefix))
		if rawURL == "" {
			log.Printf("⚠️ Missing REDIS_URL for prefix: %s, skip...", service.EnvKeyPrefix)
			continue
		}

		parsed, err := url.Parse(rawURL)
		if err != nil {
			log.Printf("Failed to parse REDIS_URL for prefix %s: %v", service.EnvKeyPrefix, err)
			continue
		}

		address := parsed.Host
		password, _ := parsed.User.Password()
		useTLS := parsed.Scheme == "rediss"

		opts := &redis.Options{
			Addr:     address,
			Password: password,
			DB:       0,
		}

		// TLS enable if needed
		if useTLS {
			opts.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		client := redis.NewClient(opts)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			log.Printf("Failed to connect Redis prefix %s: %v", service.EnvKeyPrefix, err)
			continue
		}

		redisClients[service.EnvKeyPrefix] = &RedisDynamic{
			Client:  client,
			UseTLS:  useTLS,
			Address: address,
		}

		log.Printf("[SERVICE]: ⚙️  Redis prefix: %s connected to %s", service.EnvKeyPrefix, address)
	}

	return nil
}

func GetRedisClient(prefix string) (*redis.Client, error) {
	client, ok := redisClients[prefix]
	if !ok {
		return nil, fmt.Errorf("redis client for prefix %s not found", prefix)
	}
	return client.Client, nil
}
