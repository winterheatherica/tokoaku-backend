package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedCurrentEvent(db *gorm.DB) {
	var eventType models.EventType
	if err := db.First(&eventType, "name = ?", "No Event").Error; err != nil {
		log.Println("❌ Gagal menemukan EventType 'No Event' untuk CurrentEvent:", err)
		return
	}

	current := models.CurrentEvent{
		EventTypeID: eventType.ID,
		Start:       time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := db.FirstOrCreate(&current, models.CurrentEvent{EventTypeID: current.EventTypeID}).Error; err != nil {
		log.Printf("Gagal seeding current_event: %v\n", err)
		return
	}

	log.Println("[SEEDER] ⚙️  current event seeded")
}
