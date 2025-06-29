package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

type AgentMetadata struct {
	Stomatology int    `json:"stomatology"`
	AccessToken string `json:"access_token"`
}

func (m *AgentMetadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *AgentMetadata) Scan(value interface{}) error {
	if value == nil {
		*m = AgentMetadata{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("тип должен быть []byte")
	}

	return json.Unmarshal(bytes, m)
}

type Agent struct {
	// Уникальный идентификатор агента
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Настройки API и модели
	APIKey              string        `json:"api_key" gorm:"not null;index"`
	Model               string        `json:"model" gorm:"not null;index"`
	SystemPrompt        string        `json:"system_prompt" gorm:"type:text"`
	UserPrompt          string        `json:"user_prompt" gorm:"type:text"`
	ContextSize         int           `json:"context_size" gorm:"not null;default:4096"`
	Temperature         float64       `json:"temperature" gorm:"default:0.5"`
	MaxCompletionTokens int           `json:"max_completion_tokens" gorm:"default:1000"`
	Metadata            AgentMetadata `json:"metadata" gorm:"type:jsonb"`

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
