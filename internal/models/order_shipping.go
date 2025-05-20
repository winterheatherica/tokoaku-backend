package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderShipping struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID          uint      `gorm:"not null" json:"order_id"`
	SellerID         string    `gorm:"not null" json:"seller_id"`
	ShippingOptionID uint      `gorm:"not null" json:"shipping_option_id"`
	BankAccountID    uuid.UUID `gorm:"not null" json:"bank_account_id"`
	TrackingNumber   *string   `gorm:"type:varchar(100)" json:"tracking_number"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order          Order          `gorm:"foreignKey:OrderID" json:"order"`
	Seller         User           `gorm:"foreignKey:SellerID;references:ID" json:"seller"`
	ShippingOption ShippingOption `gorm:"foreignKey:ShippingOptionID" json:"shipping_option"`
	BankAccount    BankAccount    `gorm:"foreignKey:BankAccountID" json:"bank_account"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderShippingID" json:"order_items"`
}
