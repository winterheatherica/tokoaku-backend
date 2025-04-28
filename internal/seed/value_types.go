package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedValueTypes(db *gorm.DB) {

	valueTypes := []models.ValueType{
		{ID: 1, Name: "Percentage", CreatedAt: time.Now()},
		{ID: 2, Name: "Flat", CreatedAt: time.Now()},
	}

	for _, value_type := range valueTypes {
		if err := db.FirstOrCreate(&value_type, models.ValueType{ID: value_type.ID}).Error; err != nil {
			log.Printf("Gagal seeding value_type ID %d: %v\n", value_type.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  value types seeded")
}
