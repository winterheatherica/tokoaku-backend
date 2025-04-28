package models

import "time"

type Notification struct {
	UserID             string    `gorm:"type:varchar(100);not null" json:"user_id"`
	NotificationTypeID uint      `gorm:"not null" json:"notification_type_id"`
	Message            string    `gorm:"type:text;not null" json:"message"`
	Read               bool      `gorm:"default:false" json:"read"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`

	User             User             `gorm:"foreignKey:UserID;references:ID" json:"user"`
	NotificationType NotificationType `gorm:"foreignKey:NotificationTypeID;references:ID" json:"notification_type"`
}
