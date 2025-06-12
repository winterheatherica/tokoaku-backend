package models

import "time"

type ProductVariant struct {
	ID          string    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	VariantName string    `gorm:"type:varchar(255);not null" json:"variant_name"`
	ProductID   string    `gorm:"type:varchar(100);not null;uniqueIndex:product_slug" json:"product_id"`
	Stock       uint      `gorm:"not null" json:"stock"`
	Slug        string    `gorm:"type:varchar(255);not null;uniqueIndex:product_slug" json:"slug"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`

	ProductVariantPrices []ProductPrice        `gorm:"foreignKey:ProductVariantID" json:"product_variant_prices"`
	ProductVariantImages []ProductVariantImage `gorm:"foreignKey:ProductVariantID" json:"product_variant_images"`
}
