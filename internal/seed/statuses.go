package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedStatuses(db *gorm.DB) {
	statuses := []models.Status{
		{ID: 11, StatusName: "pending", TableCategory: "order", Description: "Order dibuat, menunggu pembayaran", CreatedAt: time.Now()},
		{ID: 12, StatusName: "paid", TableCategory: "order", Description: "Pembayaran berhasil diterima", CreatedAt: time.Now()},
		{ID: 13, StatusName: "canceled", TableCategory: "order", Description: "Order dibatalkan sebelum dikirim", CreatedAt: time.Now()},
		{ID: 14, StatusName: "failed", TableCategory: "order", Description: "Pembayaran gagal", CreatedAt: time.Now()},
		{ID: 15, StatusName: "refunded", TableCategory: "order", Description: "Dana dikembalikan ke buyer", CreatedAt: time.Now()},

		{ID: 21, StatusName: "shipping", TableCategory: "shipping", Description: "Barang dalam proses pengiriman", CreatedAt: time.Now()},
		{ID: 22, StatusName: "in transit", TableCategory: "shipping", Description: "Barang masih dalam perjalanan antar hub", CreatedAt: time.Now()},
		{ID: 23, StatusName: "delivered", TableCategory: "shipping", Description: "Barang telah diterima oleh buyer", CreatedAt: time.Now()},
		{ID: 24, StatusName: "out for delivery", TableCategory: "shipping", Description: "Barang keluar dari hub terakhir menuju buyer", CreatedAt: time.Now()},
		{ID: 25, StatusName: "return_requested", TableCategory: "shipping", Description: "Buyer meminta retur barang", CreatedAt: time.Now()},
		{ID: 26, StatusName: "returned", TableCategory: "shipping", Description: "Barang dikembalikan ke seller", CreatedAt: time.Now()},
	}

	for _, status := range statuses {
		if err := db.FirstOrCreate(&status, models.Status{ID: status.ID}).Error; err != nil {
			log.Printf("Gagal seeding status ID %d: %v\n", status.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  statuses seeded")
}
