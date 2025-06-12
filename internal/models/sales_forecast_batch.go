package models

import (
	"time"
)

type SalesForecastBatch struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	Analysis  string    `gorm:"type:text" json:"analysis"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Forecasts []SalesForecast `gorm:"foreignKey:BatchID;references:ID" json:"forecasts"`
}
