package models

import "time"

type ProductVariant struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VariantName          string    `gorm:"type:varchar(255);not null" json:"variant_name"`
	ProductID            uint      `gorm:"not null" json:"product_id"`
	ProductVariantSlugID uint      `gorm:"not null" json:"product_variant_slug_id"`
	Stock                uint      `gorm:"not null" json:"stock"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductVariantSlug ProductVariantSlug `gorm:"foreignKey:ProductVariantSlugID;references:ID" json:"product_variant_slug"`
	Product            Product            `gorm:"foreignKey:ProductID" json:"product"`
}
