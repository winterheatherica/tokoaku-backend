package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/winterheatherica/tokoaku-backend/config"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
)

var DB *gorm.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DB.Host,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
		config.DB.Port,
		config.DB.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek PostgreSQL:", err)
	}

	if err := DB.AutoMigrate(
		&models.ActiveAddress{},
		&models.ActiveBankAccount{},
		&models.Address{},
		&models.BankAccount{},
		&models.BankList{},
		&models.BankTransferFee{},
		&models.Cart{},
		&models.CategoryDiscount{},
		&models.Category{},
		&models.CloudService{},
		&models.CurrentEvent{},
		&models.DailyCategoryPrice{},
		&models.DefaultFee{},
		&models.DiscountSponsor{},
		&models.Discount{},
		&models.EventType{},
		&models.NotificationType{},
		&models.Notification{},
		&models.OrderItem{},
		&models.OrderLog{},
		&models.OrderPromo{},
		&models.OrderShippingStatus{},
		&models.OrderShipping{},
		&models.Order{},
		&models.PaymentMethod{},
		&models.PendingUser{},
		&models.ProductCategory{},
		&models.ProductForm{},
		&models.ProductPrice{},
		&models.ProductTypeDiscount{},
		&models.ProductType{},
		&models.ProductVariantDiscount{},
		&models.ProductVariantImage{},
		&models.ProductVariant{},
		&models.Product{},
		&models.Promo{},
		&models.Provider{},
		&models.Review{},
		&models.Role{},
		&models.SellerShippingOption{},
		&models.Sentiment{},
		&models.ShippingOption{},
		&models.Status{},
		&models.SummarizationDetail{},
		&models.Summarization{},
		&models.UserPromo{},
		&models.User{},
		&models.ValueType{},
	); err != nil {
		log.Fatal("Migration error:", err)
	}

	log.Println("[SERVICE]: ⚙️  PostgreSQL connected")
}
