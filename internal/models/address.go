package models

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      string    `gorm:"type:varchar(100);not null" json:"user_id"`
	AddressLine string    `gorm:"type:text;not null" json:"address_line"`
	City        string    `gorm:"not null" json:"city"`
	Province    string    `gorm:"not null" json:"province"`
	PostalCode  string    `gorm:"not null" json:"postal_code"`
	Latitude    float64   `gorm:"type:float" json:"latitude"`
	Longitude   float64   `gorm:"type:float" json:"longitude"`
	IsActive    bool      `gorm:"default:false" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"user"`
}
