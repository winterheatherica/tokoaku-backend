package config

import "log"

func LoadAll() {

	LoadAppConfig()
	LoadFirebaseConfig()
	LoadDatabaseConfig()
	LoadRedisConfig()
	LoadResendConfig()
	LoadCloudinaryConfig()

	log.Println("[CONFIG]: ✅ All config initialized")
}
