package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID      string    `gorm:"not null" json:"customer_id"`
	PaymentMethodID uint      `gorm:"not null" json:"payment_method_id"`
	AddressID       uuid.UUID `gorm:"not null" json:"address_id"`
	TotalPrice      uint      `gorm:"not null;default:0" json:"total_price"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	Customer      User          `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID" json:"payment_method"`
	Address       Address       `gorm:"foreignKey:AddressID;references:ID" json:"address"`
}
