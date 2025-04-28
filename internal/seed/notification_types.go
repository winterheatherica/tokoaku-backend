package seed

import (
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedNotificationTypes(db *gorm.DB) error {
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

	if err := db.Create(&notificationTypes).Error; err != nil {
		return err
	}
	return nil
}
