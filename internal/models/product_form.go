package models

import "time"

type ProductForm struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Form         string    `gorm:"type:varchar(100);not null" json:"form"`
	Slug         string    `gorm:"uniqueIndex;not null" json:"slug"`
	ImageURL     string    `gorm:"type:text;not null" json:"image_url"`
	CloudImageID uint      `gorm:"not null" json:"cloud_image_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	CloudService CloudService `gorm:"foreignKey:CloudImageID" json:"cloud_service"`
}
