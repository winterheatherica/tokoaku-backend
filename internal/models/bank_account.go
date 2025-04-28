package models

import "time"

type BankAccount struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	BankID        uint      `gorm:"not null" json:"bank_id"`
	AccountNumber string    `gorm:"type:varchar(50);not null" json:"account_number"`
	AccountName   string    `gorm:"type:varchar(100);not null" json:"account_name"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	User User     `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Bank BankList `gorm:"foreignKey:BankID;references:ID" json:"bank"`
}
