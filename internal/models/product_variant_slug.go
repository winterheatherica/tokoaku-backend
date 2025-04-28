package models

import "time"

type ProductVariantSlug struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
