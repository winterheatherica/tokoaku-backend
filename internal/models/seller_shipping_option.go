package models

import "time"

type SellerShippingOption struct {
	ShippingOptionID uint      `gorm:"not null" json:"shipping_option_id"`
	SellerID         uint      `gorm:"not null" json:"seller_id"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	ShippingOption ShippingOption `gorm:"foreignKey:ShippingOptionID;references:ID" json:"shipping_option"`
	Seller         User           `gorm:"foreignKey:SellerID;references:ID" json:"seller"`
}
