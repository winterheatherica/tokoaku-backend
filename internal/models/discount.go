package models

import "time"

type Discount struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	ValueTypeID uint      `gorm:"not null" json:"value_type_id"`
	Value       uint      `gorm:"not null" json:"value"`
	SponsorID   uint      `gorm:"not null" json:"sponsor_id"`
	StartAt     time.Time `gorm:"column:start_at"`
	EndAt       time.Time `gorm:"column:end_at"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	ValueType       ValueType       `gorm:"foreignKey:ValueTypeID" json:"value_type"`
	DiscountSponsor DiscountSponsor `gorm:"foreignKey:SponsorID" json:"discount_sponsor"`
}
