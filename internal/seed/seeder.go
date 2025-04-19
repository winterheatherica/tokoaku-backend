package seed

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) {
	SeedRoles(db)
	SeedProviders(db)
	log.Println("[SEEDER] ✅ All seeders executed")
}
