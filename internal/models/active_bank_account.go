package models

import "time"

type ActiveBankAccount struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	BankAccount BankAccount `gorm:"foreignKey:ID;references:ID" json:"bank_account"`
}
