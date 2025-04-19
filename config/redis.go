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
		Redis.Host = ""
		Redis.Password = ""
		Redis.UseTLS = false
		return
	}

	Redis.Host = parsed.Host
	if parsed.User != nil {
		Redis.Password, _ = parsed.User.Password()
	}

	Redis.UseTLS = parsed.Scheme == "rediss"

	log.Println("[CONFIG]: ⚙️  Redis config initialized")
}
