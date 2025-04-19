package services

import (
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/firebase"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
)

func InitAll() {

	firebase.InitFirebase()
	database.Connect()
	redis.Connect()

	log.Println("[SERVICE]: ✅ All service connected")
}
