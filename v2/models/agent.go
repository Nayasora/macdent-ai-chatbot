package models

import (
	"github.com/google/uuid"
	"time"
)

type Agent struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	APIKey       string      `json:"api_key" gorm:"not null"`
	Model        string      `json:"model" gorm:"not null"`
	SystemPrompt string      `json:"system_prompt" gorm:"type:text"`
	UserPrompt   string      `json:"user_prompt" gorm:"type:text"`
	ContextSize  int         `json:"context_size" gorm:"not null;default:4096"`
	Temperature  float32     `json:"temperature" gorm:"default:0.7"`
	TopP         float32     `json:"top_p" gorm:"default:1.0"`
	MaxTokens    int         `json:"max_tokens" gorm:"default:1000"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Permissions  *Permission `json:"permissions" gorm:"foreignKey:AgentID;constraint:OnDelete:CASCADE"`
}
