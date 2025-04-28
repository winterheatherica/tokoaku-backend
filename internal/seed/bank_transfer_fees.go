package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedBankTransferFees(db *gorm.DB) {
	transferFees := []models.BankTransferFee{
		{FromBankID: 4, ToBankID: 7, FeeID: 3, CreatedAt: time.Now()},
		{FromBankID: 4, ToBankID: 1, FeeID: 5, CreatedAt: time.Now()},
		{FromBankID: 4, ToBankID: 3, FeeID: 2, CreatedAt: time.Now()},
		{FromBankID: 4, ToBankID: 4, FeeID: 6, CreatedAt: time.Now()},
		{FromBankID: 4, ToBankID: 52, FeeID: 4, CreatedAt: time.Now()},
		{FromBankID: 4, ToBankID: 53, FeeID: 1, CreatedAt: time.Now()},
	}

	for _, f := range transferFees {
		if err := db.FirstOrCreate(&f, models.BankTransferFee{
			FromBankID: f.FromBankID,
			ToBankID:   f.ToBankID,
		}).Error; err != nil {
			log.Printf("Gagal seeding transfer_fee %d ➝ %d: %v\n", f.FromBankID, f.ToBankID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  bank transfer fees seeded")
}
