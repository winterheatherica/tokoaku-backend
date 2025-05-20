package models

import "time"

type ProductPrice struct {
	ProductVariantID string    `gorm:"not null;index:idx_price_variant" json:"product_variant_id"`
	Price            uint      `gorm:"not null" json:"price"`
	CreatedAt        time.Time `gorm:"autoCreateTime;index:idx_price_created" json:"created_at"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
}
