package models

import (
	"time"
)

type OrderShippingStatus struct {
	OrderShippingID uint      `gorm:"not null" json:"order_shipping_id"`
	StatusID        uint      `gorm:"not null" json:"status_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	OrderShipping OrderShipping `gorm:"foreignKey:OrderShippingID;references:ID" json:"order_shipping"`
	Status        Status        `gorm:"foreignKey:StatusID;references:ID" json:"status"`
}
