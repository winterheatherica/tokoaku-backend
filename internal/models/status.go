package models

import (
	"time"
)

type Status struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StatusName    string    `gorm:"type:varchar(100);not null" json:"status_name"`
	TableCategory string    `gorm:"type:varchar(50);not null" json:"table_category"`
	Description   string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
