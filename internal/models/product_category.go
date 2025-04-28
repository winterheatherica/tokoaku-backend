package models

import "time"

type ProductCategory struct {
	ProductID  uint      `gorm:"not null" json:"product_id"`
	CategoryID uint      `gorm:"not null" json:"category_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	Product  Product  `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Category Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
}
