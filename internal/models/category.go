package models

import "time"

type Category struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug            string    `gorm:"uniqueIndex;not null" json:"slug"`
	Code            *string   `gorm:"size:50" json:"code,omitempty"`
	CategoryLabelID uint      `gorm:"not null" json:"category_label_id"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	CategoryLabel CategoryLabel `gorm:"foreignKey:CategoryLabelID;references:ID" json:"category_label"`
}
