package models

import "time"

type DiscountSponsor struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}
