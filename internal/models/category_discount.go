package models

import "time"

type CategoryDiscount struct {
	CategoryID uint      `gorm:"not null" json:"category_id"`
	DiscountID uint      `gorm:"not null" json:"discount_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	Category Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	Discount Discount `gorm:"foreignKey:DiscountID;references:ID" json:"discount"`
}
