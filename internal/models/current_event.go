package models

import "time"

type CurrentEvent struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventTypeID uint      `gorm:"not null" json:"event_type_id"`
	Start       time.Time `gorm:"not null" json:"start"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	EventType EventType `gorm:"foreignKey:EventTypeID;references:ID" json:"event_type"`
}
