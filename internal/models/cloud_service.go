package models

import "time"

type CloudService struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	ProviderID    uint      `gorm:"not null" json:"provider_id"`
	EnvKeyPrefix  string    `gorm:"unique;not null" json:"env_key_prefix"`
	UsageFor      string    `gorm:"not null" json:"usage_for"`
	StorageUsage  uint      `json:"storage_usage"`
	LastCheckedAt time.Time `json:"last_checked_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	Provider Provider `gorm:"foreignKey:ProviderID;references:ID" json:"provider"`
}
