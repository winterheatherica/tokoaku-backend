package models

import "time"

type Promo struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Code          string    `gorm:"type:varchar(100);unique;not null" json:"code"`
	Description   string    `gorm:"type:text" json:"description"`
	ValueTypeID   uint      `gorm:"not null" json:"value_type_id"`
	Value         uint      `gorm:"not null" json:"value"`
	MinPriceValue uint      `gorm:"default:0" json:"min_price_value"`
	MaxValue      uint      `gorm:"default:0" json:"max_value"`
	StartAt       time.Time `gorm:"not null" json:"start_at"`
	EndAt         time.Time `gorm:"not null" json:"end_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	ValueType ValueType `gorm:"foreignKey:ValueTypeID;references:ID" json:"value_type"`
}
