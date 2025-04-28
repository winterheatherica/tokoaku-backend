package models

import "time"

type BankTransferFee struct {
	FromBankID uint      `gorm:"not null" json:"from_bank_id"`
	ToBankID   uint      `gorm:"not null" json:"to_bank_id"`
	FeeID      uint      `gorm:"not null" json:"fee_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	FromBank BankList   `gorm:"foreignKey:FromBankID;references:ID" json:"from_bank"`
	ToBank   BankList   `gorm:"foreignKey:ToBankID;references:ID" json:"to_bank"`
	Fee      DefaultFee `gorm:"foreignKey:FeeID;references:ID" json:"default_fee"`
}
