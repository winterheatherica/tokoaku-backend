package models

import (
	"time"
)

type HighlightedProduct struct {
	ProductID string    `gorm:"not null" json:"product_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}
