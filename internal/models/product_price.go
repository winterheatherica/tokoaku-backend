package models

import "time"

type ProductPrice struct {
	ProductVariantID uint      `gorm:"not null" json:"product_variant_id"`
	Price            uint      `gorm:"not null" json:"price"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
}
