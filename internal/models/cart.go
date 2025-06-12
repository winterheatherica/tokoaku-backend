package models

import "time"

type Cart struct {
	CustomerID       string    `gorm:"not null;index:idx_carts_customer_id;uniqueIndex:idx_unique_customer_variant,priority:1" json:"customer_id"`
	ProductVariantID string    `gorm:"not null;index:idx_carts_variant_id;uniqueIndex:idx_unique_customer_variant,priority:2" json:"product_variant_id"`
	Quantity         uint      `gorm:"not null" json:"quantity"`
	IsSelected       bool      `gorm:"default:false;index:idx_carts_selected_customer,priority:2;index:idx_carts_conversion_check,priority:2" json:"is_selected"`
	IsConverted      bool      `gorm:"default:false;index:idx_carts_conversion_check,priority:3" json:"is_converted"`
	CreatedAt        time.Time `gorm:"autoCreateTime;uniqueIndex:idx_unique_customer_variant,priority:3"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
	Customer       User           `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}
