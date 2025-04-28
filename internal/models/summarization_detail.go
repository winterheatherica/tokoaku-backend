package models

import "time"

type SummarizationDetail struct {
	SummarizationID uint      `gorm:"not null" json:"summarization_id"`
	Text            string    `gorm:"type:text;not null" json:"text"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	Summarization Summarization `gorm:"foreignKey:SummarizationID;references:ID" json:"summarization"`
}
