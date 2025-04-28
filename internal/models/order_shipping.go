package models

import (
	"time"
)

type OrderShipping struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID          uint      `gorm:"not null" json:"order_id"`
	ShippingOptionID uint      `gorm:"not null" json:"shipping_option_id"`
	SellerID         uint      `gorm:"not null" json:"seller_id"`
	TrackingNumber   string    `gorm:"type:varchar(100)" json:"tracking_number"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order          Order          `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	ShippingOption ShippingOption `gorm:"foreignKey:ShippingOptionID;references:ID" json:"shipping_option"`
	Seller         User           `gorm:"foreignKey:SellerID;references:ID" json:"seller"`
}
