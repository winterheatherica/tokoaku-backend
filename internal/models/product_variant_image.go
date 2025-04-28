package models

import (
	"time"
)

type ProductVariantImage struct {
	ProductVariantID string    `gorm:"not null" json:"product_variant_id"`
	ImageURL         string    `gorm:"type:text;not null" json:"image_url"`
	CloudImageID     uint      `gorm:"not null" json:"cloud_image_id"`
	IsCover          bool      `gorm:"default:false" json:"is_variant_cover"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant"`
	CloudService   CloudService   `gorm:"foreignKey:CloudImageID" json:"cloud_service"`
}
