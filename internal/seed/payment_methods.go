package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedPaymentMethods(db *gorm.DB) {

	paymentMethods := []models.PaymentMethod{
		{ID: 1, Name: "Virtual Account", CreatedAt: time.Now()},
	}

	for _, paymentMethod := range paymentMethods {
		if err := db.FirstOrCreate(&paymentMethod, models.PaymentMethod{ID: paymentMethod.ID}).Error; err != nil {
			log.Printf("Gagal seeding payment method ID %d: %v\n", paymentMethod.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  payment methods seeded")
}
