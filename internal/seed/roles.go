package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/config"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {

	roles := []models.Role{
		{ID: 1, Name: "Customer", CreatedAt: time.Now()},
		{ID: 2, Name: "Seller", CreatedAt: time.Now()},
		{ID: 3, Name: "Admin", CreatedAt: time.Now()},
		{ID: 4, Name: config.App.PlatformName, CreatedAt: time.Now()},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{ID: role.ID}).Error; err != nil {
			log.Printf("Gagal seeding role ID %d: %v\n", role.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  roles seeded")
}
