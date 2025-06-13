package models

import "time"

type Discount struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `gorm:"type:text;not null" json:"description"`
	ValueTypeID   uint      `gorm:"not null" json:"value_type_id"`
	Value         uint      `gorm:"not null" json:"value"`
	SponsorID     uint      `gorm:"not null" json:"sponsor_id"`
	StartAt       time.Time `gorm:"column:start_at" json:"start_at"`
	EndAt         time.Time `gorm:"column:end_at" json:"end_at"`
	Slug          string    `gorm:"uniqueIndex;not null" json:"slug"`
	ImageCoverURL string    `gorm:"type:text" json:"image_cover_url"`
	CloudImageID  uint      `gorm:"not null" json:"cloud_image_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	ValueType       ValueType       `gorm:"foreignKey:ValueTypeID" json:"value_type"`
	DiscountSponsor DiscountSponsor `gorm:"foreignKey:SponsorID" json:"discount_sponsor"`
	CloudService    CloudService    `gorm:"foreignKey:CloudImageID" json:"cloud_service"`
}
