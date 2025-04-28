package models

import (
	"time"

	"github.com/google/uuid"
)

type ActiveAddress struct {
	AddressID uuid.UUID `gorm:"type:char(36);primaryKey" json:"address_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Address Address `gorm:"foreignKey:AddressID;references:ID" json:"address"`
}
