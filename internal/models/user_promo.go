package models

import "time"

type UserPromo struct {
	PromoID    uint      `gorm:"not null" json:"promo_id"`
	CustomerID uint      `gorm:"not null" json:"customer_id"`
	Redeemed   bool      `gorm:"default:false" json:"redeemed"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	Promo    Promo `gorm:"foreignKey:PromoID;references:ID" json:"promo"`
	Customer User  `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
}
