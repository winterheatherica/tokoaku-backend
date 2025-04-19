package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/winterheatherica/tokoaku-backend/config"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
)

var DB *gorm.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DB.Host,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
		config.DB.Port,
		config.DB.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek PostgreSQL:", err)
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.PendingUser{},
		&models.User{},
		&models.Role{},
	); err != nil {
		log.Fatal("Migration error:", err)
	}

	log.Println("[SERVICE]: ⚙️  PostgreSQL connected")
}
