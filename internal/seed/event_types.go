package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedEventTypes(db *gorm.DB) {
	eventTypes := []models.EventType{
		{Name: "No Event", DiscountLimit: 1, CreatedAt: time.Now()},
		{Name: "Big Event", DiscountLimit: 2, CreatedAt: time.Now()},
		{Name: "Flash Event", DiscountLimit: 3, CreatedAt: time.Now()},
	}

	for _, e := range eventTypes {
		event := e
		if err := db.FirstOrCreate(&event, models.EventType{Name: event.Name}).Error; err != nil {
			log.Printf("Gagal seeding event_type %s: %v\n", event.Name, err)
		}
	}

	log.Println("[SEEDER] ⚙️  event types seeded")
}
