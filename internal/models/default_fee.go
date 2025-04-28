package models

import "time"

type DefaultFee struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceName string    `gorm:"not null" json:"service_name"`
	Fee         uint      `gorm:"not null" json:"fee"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
