package models

import "time"

type Order struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID      uint      `gorm:"not null" json:"customer_id"`
	PaymentMethodID uint      `gorm:"not null" json:"payment_method_id"`
	TotalPrice      uint      `gorm:"not null;default:0" json:"total_price"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	Customer      User          `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID" json:"payment_method"`
}
