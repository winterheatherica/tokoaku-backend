package models

import "time"

type OrderPromo struct {
	OrderID   uint      `gorm:"not null" json:"order_id"`
	PromoID   uint      `gorm:"not null" json:"promo_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order Order `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	Promo Promo `gorm:"foreignKey:PromoID;references:ID" json:"promo"`
}
