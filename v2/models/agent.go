package models

import (
	"github.com/google/uuid"
	"time"
)

// Agent представляет AI-агента с его конфигурацией
type Agent struct {
	// Уникальный идентификатор агента
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Настройки API и модели
	APIKey       string  `json:"api_key" gorm:"not null;index"`
	Model        string  `json:"model" gorm:"not null;index"`
	SystemPrompt string  `json:"system_prompt" gorm:"type:text"`
	UserPrompt   string  `json:"user_prompt" gorm:"type:text"`
	ContextSize  int     `json:"context_size" gorm:"not null;default:4096"`
	Temperature  float32 `json:"temperature" gorm:"default:0.7"`
	TopP         float32 `json:"top_p" gorm:"default:1.0"`
	MaxTokens    int     `json:"max_tokens" gorm:"default:1000"`

	// Метаданные
	CreatedAt time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	// Связи с другими таблицами
	Permission       Permission        `json:"permission" gorm:"foreignKey:AgentID;references:ID;constraint:OnDelete:CASCADE"`
	KnowledgePrompts []KnowledgePrompt `json:"knowledge_prompts,omitempty" gorm:"foreignKey:AgentID;references:ID"`
	KnowledgeFiles   []KnowledgeFile   `json:"knowledge_files,omitempty" gorm:"foreignKey:AgentID;references:ID"`
	Dialogs          []Dialog          `json:"dialogs,omitempty" gorm:"foreignKey:AgentID;references:ID"`
}
