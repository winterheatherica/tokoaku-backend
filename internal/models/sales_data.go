package models

import (
	"time"
)

type SalesData struct {
	Date       time.Time `gorm:"primaryKey;type:date" json:"date"`
	TotalSales int64     `gorm:"not null" json:"total_sales"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
