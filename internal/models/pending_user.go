package models

import "time"

type PendingUser struct {
	Email        string    `gorm:"primaryKey;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"type:text" json:"password_hash"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
