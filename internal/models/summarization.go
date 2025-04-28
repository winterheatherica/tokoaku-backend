package models

import "time"

type Summarization struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint      `gorm:"not null" json:"product_id"`
	SentimentID uint      `gorm:"not null" json:"sentiment_id"`
	CeviewCount uint      `gorm:"default:0" json:"review_count"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Sentiment Sentiment `gorm:"foreignKey:SentimentID;references:ID" json:"sentiment"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}
