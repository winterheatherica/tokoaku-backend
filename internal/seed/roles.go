package seed

import (
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {

	roles := []models.Role{
		{ID: 1, Name: "user"},
		{ID: 2, Name: "vendor"},
		{ID: 3, Name: "admin"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{ID: role.ID}).Error; err != nil {
			log.Printf("Gagal seeding role ID %d: %v\n", role.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  Role seeded")
}
