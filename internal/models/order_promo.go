package models

import "time"

type OrderPromo struct {
	OrderID          uint      `gorm:"not null" json:"order_id"`
	ProductVariantID uint      `gorm:"not null" json:"product_variant_id"`
	PromoID          uint      `gorm:"not null" json:"promo_id"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order          Order          `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
	Promo          Promo          `gorm:"foreignKey:PromoID;references:ID" json:"promo"`
}
