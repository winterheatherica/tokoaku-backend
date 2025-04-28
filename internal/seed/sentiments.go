package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedSentiments(db *gorm.DB) {

	sentiments := []models.Sentiment{
		{ID: 1, Name: "Positive", CreatedAt: time.Now()},
		{ID: 2, Name: "Negative", CreatedAt: time.Now()},
	}

	for _, sentiment := range sentiments {
		if err := db.FirstOrCreate(&sentiment, models.Sentiment{ID: sentiment.ID}).Error; err != nil {
			log.Printf("Gagal seeding sentiment ID %d: %v\n", sentiment.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  sentiments seeded")
}
