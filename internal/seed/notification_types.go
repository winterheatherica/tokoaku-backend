package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedNotificationTypes(db *gorm.DB) {
	notificationTypes := []models.NotificationType{
		{Name: "Order Update", CreatedAt: time.Now()},
		{Name: "Payment Reminder", CreatedAt: time.Now()},
		{Name: "Promo", CreatedAt: time.Now()},
		{Name: "Message", CreatedAt: time.Now()},
		{Name: "System Update", CreatedAt: time.Now()},
		{Name: "Shipping Update", CreatedAt: time.Now()},
		{Name: "Refund Processed", CreatedAt: time.Now()},
		{Name: "Review Reminder", CreatedAt: time.Now()},
		{Name: "Account Alert", CreatedAt: time.Now()},
		{Name: "Loyalty Reward", CreatedAt: time.Now()},
	}

	for _, nt := range notificationTypes {
		ntItem := nt
		if err := db.FirstOrCreate(&ntItem, models.NotificationType{Name: ntItem.Name}).Error; err != nil {
			log.Printf("❌ Gagal seeding notification_type %s: %v\n", ntItem.Name, err)
		}
	}

	log.Println("[SEEDER] 🔔 Notification types seeded")
}
