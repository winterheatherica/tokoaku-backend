package models

import "time"

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(100)" json:"id,omitempty"`
	Username     *string   `gorm:"uniqueIndex" json:"username,omitempty"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Phone        *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	PasswordHash *string   `gorm:"type:text" json:"password_hash,omitempty"`
	ProviderID   uint      `gorm:"not null"`
	RoleID       uint      `gorm:"default:1" json:"role_id"`
	Name         *string   `gorm:"type:varchar(100)" json:"name,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	Provider Provider `gorm:"foreignKey:ProviderID;references:ID" json:"provider"`
	Role     Role     `gorm:"foreignKey:RoleID;references:ID" json:"role"`
}
