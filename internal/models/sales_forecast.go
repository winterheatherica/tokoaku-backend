package models

import (
	"time"
)

type SalesForecast struct {
	Date           time.Time `gorm:"primaryKey;type:date" json:"date"`
	PredictedSales int64     `gorm:"not null" json:"predicted_sales"`
	BatchID        string    `gorm:"not null" json:"batch_id"`

	Batch SalesForecastBatch `gorm:"foreignKey:BatchID;references:ID" json:"batch"`
}
