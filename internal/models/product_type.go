package models

import "time"

type ProductType struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Slug         string    `gorm:"uniqueIndex;not null" json:"slug"`
	ImageURL     string    `gorm:"type:text;not null" json:"image_url"`
	CloudImageID uint      `gorm:"not null" json:"cloud_image_id"`
	ValueTypeID  uint      `gorm:"not null" json:"value_type_id"`
	Value        uint      `gorm:"not null" json:"value"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	ValueType    ValueType    `gorm:"foreignKey:ValueTypeID" json:"value_type"`
	CloudService CloudService `gorm:"foreignKey:CloudImageID" json:"cloud_service"`
}
