package models

import "time"

type OrderLog struct {
	OrderID   uint      `gorm:"not null" json:"order_id"`
	StatusID  uint      `gorm:"not null" json:"status_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order  Order  `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	Status Status `gorm:"foreignKey:StatusID;references:ID" json:"status"`
}
