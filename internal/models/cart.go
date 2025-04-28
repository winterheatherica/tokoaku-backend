package models

import "time"

type Cart struct {
	CustomerID       uint      `gorm:"not null" json:"customer_id"`
	ProductVariantID uint      `gorm:"not null" json:"product_variant_id"`
	Quantity         uint      `gorm:"not null" json:"quantity"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
	Customer       User           `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}
