package config

import "log"

func LoadAll() {

	LoadAppConfig()
	LoadFirebaseConfig()
	LoadDatabaseConfig()
	LoadResendConfig()

	log.Println("[CONFIG]: ✅ All config initialized")
}
