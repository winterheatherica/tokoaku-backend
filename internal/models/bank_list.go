package models

import "time"

type BankList struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	FullName  string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Code      string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
