package models

import "time"

type ProductVariantDiscount struct {
	ProductVariantID uint      `gorm:"not null" json:"product_variant_id"`
	DiscountID       uint      `gorm:"not null" json:"discount_id"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant"`
	Discount       Discount       `gorm:"foreignKey:DiscountID" json:"discount"`
}
