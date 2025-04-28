package models

import "time"

type ShippingOption struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourierName        string    `gorm:"type:varchar(100);not null" json:"courier_name"`
	CourierServiceName string    `gorm:"type:varchar(100);not null" json:"courier_service_name"`
	Fee                float64   `gorm:"not null" json:"fee"`
	EstimatedTime      string    `gorm:"type:varchar(100);not null" json:"estimated_time"`
	ServiceType        string    `gorm:"type:varchar(50);not null" json:"service_type"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}
