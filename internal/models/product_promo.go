package models

import "time"

type ProductPromo struct {
	PromoID   uint      `gorm:"not null" json:"promo_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Promo   Promo   `gorm:"foreignKey:PromoID;references:ID" json:"promo"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}
