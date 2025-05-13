package models

import "time"

type DailyCategoryPrice struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Date           time.Time `gorm:"not null;index" json:"date"`
	IsHoliday      bool      `gorm:"not null" json:"is_holiday"`
	IsPayday       bool      `gorm:"not null" json:"is_payday"`
	IsBeautifulDay bool      `gorm:"not null" json:"is_beautiful_day"`
	CategoryID     uint      `gorm:"not null;index" json:"category_id"`
	Price          float64   `gorm:"not null" json:"price"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	Category Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
}
