package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDynamic struct {
	Client  *redis.Client
	UseTLS  bool
	Address string
}

var (
	redisClients = make(map[string]*RedisDynamic)
	once         sync.Once
)

type redisConfig struct {
	Address  string
	Password string
	UseTLS   bool
}

func parseRedisURL(rawURL string) (*redisConfig, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	password, _ := parsed.User.Password()
	useTLS := parsed.Scheme == "rediss"

	return &redisConfig{
		Address:  parsed.Host,
		Password: password,
		UseTLS:   useTLS,
	}, nil
}

func newRedisClient(cfg *redisConfig) *redis.Client {
	opts := &redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	}
	if cfg.UseTLS {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return redis.NewClient(opts)
}

func LoadRedisClient(prefix string) error {

	envKey := fmt.Sprintf("%s_REDIS_URL", prefix)
	rawURL := os.Getenv(envKey)
	if rawURL == "" {
		return fmt.Errorf("missing Redis URL for prefix: %s", prefix)
	}

	cfg, err := parseRedisURL(rawURL)
	if err != nil {
		return fmt.Errorf("invalid Redis URL for prefix %s: %v", prefix, err)
	}

	client := newRedisClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis ping failed for %s: %v", prefix, err)
	}

	redisClients[prefix] = &RedisDynamic{
		Client:  client,
		UseTLS:  cfg.UseTLS,
		Address: cfg.Address,
	}

	log.Printf("[CONFIG] ✅ Redis connected | prefix: %s | host: %s", prefix, cfg.Address)
	return nil
}

func GetRedisClient(prefix string) (*redis.Client, error) {
	cfg, ok := redisClients[prefix]
	if !ok {
		return nil, fmt.Errorf("redis client for prefix %s not found", prefix)
	}
	return cfg.Client, nil
}
