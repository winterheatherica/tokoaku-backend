package models

import "time"

type EventType struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	DiscountLimit uint      `gorm:"not null" json:"discount_limit"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
