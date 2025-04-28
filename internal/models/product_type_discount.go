package models

import "time"

type ProductTypeDiscount struct {
	ProductTypeID uint      `gorm:"primaryKey" json:"product_type_id"`
	DiscountID    uint      `gorm:"primaryKey" json:"discount_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductType ProductType `gorm:"foreignKey:ProductTypeID;references:ID" json:"product_type"`
	Discount    Discount    `gorm:"foreignKey:DiscountID;references:ID" json:"discount"`
}
