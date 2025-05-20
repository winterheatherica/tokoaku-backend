package models

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UserID        string    `gorm:"not null" json:"user_id"`
	BankID        uint      `gorm:"not null" json:"bank_id"`
	AccountNumber string    `gorm:"type:varchar(50);not null" json:"account_number"`
	AccountName   string    `gorm:"type:varchar(100);not null" json:"account_name"`
	IsActive      bool      `gorm:"default:false" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	User User     `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Bank BankList `gorm:"foreignKey:BankID;references:ID" json:"bank"`
}
