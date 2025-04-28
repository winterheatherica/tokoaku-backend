package config

import (
	"log"
	"net/url"
	"os"
)

var Redis struct {
	Host     string
	Password string
	UseTLS   bool
}

func LoadRedisConfig() {
	rawURL := os.Getenv("REDIS_URL")
	parsed, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("[CONFIG]: Failed to parse REDIS_URL: %v", err)
	}

	Redis.Host = parsed.Host

	if parsed.User != nil {
		Redis.Password, _ = parsed.User.Password()
	}

	Redis.UseTLS = parsed.Scheme == "rediss"

	log.Printf("[CONFIG]: ⚙️  Redis config initialized | Host: %s | TLS: %v", Redis.Host, Redis.UseTLS)
}
