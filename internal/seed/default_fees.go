package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedDefaultFees(db *gorm.DB) {
	defaultFees := []models.DefaultFee{
		{ServiceName: "BI-Fast", Fee: 2500, Description: "Layanan BI-Fast, maksimal Rp250 juta, 24/7 realtime.", CreatedAt: time.Now()},
		{ServiceName: "RTO", Fee: 6500, Description: "Real Time Online, standar transfer antar bank.", CreatedAt: time.Now()},
		{ServiceName: "SKNBI", Fee: 2900, Description: "Sistem Kliring Nasional BI, cocok untuk transfer besar.", CreatedAt: time.Now()},
		{ServiceName: "RTGS", Fee: 25000, Description: "Real-Time Gross Settlement untuk transfer ≥ Rp100 juta.", CreatedAt: time.Now()},
		{ServiceName: "ATM Himbara", Fee: 4000, Description: "Transfer antar bank anggota Himbara melalui ATM.", CreatedAt: time.Now()},
		{ServiceName: "Internal", Fee: 0, Description: "Transfer antar rekening dalam bank yang sama.", CreatedAt: time.Now()},
	}

	for _, f := range defaultFees {
		fee := f
		if err := db.FirstOrCreate(&fee, models.DefaultFee{ServiceName: fee.ServiceName}).Error; err != nil {
			log.Printf("Gagal seeding default_fee %s: %v\n", fee.ServiceName, err)
		}
	}

	log.Println("[SEEDER] ⚙️  default fees seeded")
}
