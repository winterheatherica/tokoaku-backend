package models

import "time"

type ProductForm struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Form      string    `gorm:"type:varchar(100);not null" json:"form"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
