package models

import "time"

type Review struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductVariantID uint      `gorm:"not null" json:"product_variant_id"`
	CustomerID       uint      `gorm:"not null" json:"customer_id"`
	Text             string    `gorm:"type:varchar(255);not null" json:"text"`
	Rating           uint      `gorm:"not null" json:"rating"`
	SentimentID      *uint     `gorm:"column:sentiment_id" json:"sentiment_id,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`

	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
	Customer       User           `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	Sentiment      Sentiment      `gorm:"foreignKey:SentimentID;references:ID" json:"sentiment"`
}
