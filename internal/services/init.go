package services

import (
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/services/cloudinary"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/firebase"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
)

func InitAll() {
	firebase.InitFirebase()
	database.Connect()

	if err := redis.LoadAllRedisConfigs(database.DB); err != nil {
		log.Fatalf("Gagal load semua Redis configs: %v", err)
	}

	if err := cloudinary.LoadAllCloudinaryConfigs(database.DB); err != nil {
		log.Fatalf("Gagal load semua Cloudinary configs: %v", err)
	}

	log.Println("[SERVICE]: ✅ All service connected")
}
