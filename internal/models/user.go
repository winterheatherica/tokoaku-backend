package models

import "time"

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(100)" json:"id,omitempty"`
	Username     *string   `gorm:"uniqueIndex" json:"username,omitempty"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Phone        *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	PasswordHash *string   `gorm:"type:text" json:"password_hash,omitempty"`
	Provider     *string   `gorm:"type:varchar(50)" json:"provider,omitempty"`
	Role         int       `gorm:"default:0" json:"role"`
	Name         *string   `gorm:"type:varchar(100)" json:"name,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
