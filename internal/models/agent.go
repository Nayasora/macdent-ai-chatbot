package models

import (
	"github.com/google/uuid"
	"time"
)

type Agent struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string    `json:"name" gorm:"not null"`
	Description     string    `json:"description" gorm:"type:text"`
	APIKey          string    `json:"api_key" gorm:"not null"`
	Model           string    `json:"model" gorm:"not null"`
	ContextSize     int       `json:"context_size" gorm:"not null;default:4096"`
	SystemPrompt    string    `json:"system_prompt" gorm:"type:text"`
	AssistantPrompt string    `json:"assistant_prompt" gorm:"type:text"`
	Temperature     float32   `json:"temperature" gorm:"default:0.7"`
	MaxTokens       int       `json:"max_tokens" gorm:"default:1000"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Dialog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AgentID   uuid.UUID `json:"agent_id" gorm:"type:uuid;not null;index"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Response  string    `json:"response" gorm:"type:text"`
	Role      string    `json:"role" gorm:"not null"`
	Metadata  string    `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at"`
	Agent     Agent     `json:"-" gorm:"foreignKey:AgentID"`
}

type KnowledgeFile struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	AgentID        uuid.UUID  `json:"agent_id" gorm:"type:uuid;not null;index"`
	FileName       string     `json:"file_name" gorm:"not null"`
	OriginalName   string     `json:"original_name" gorm:"not null"`
	FileSize       int64      `json:"file_size"`
	FileType       string     `json:"file_type"`
	FilePath       string     `json:"file_path"`
	CollectionName string     `json:"collection_name" gorm:"not null"`
	ChunkCount     int        `json:"chunk_count" gorm:"default:0"`
	Status         string     `json:"status" gorm:"default:processing"`
	ErrorMessage   string     `json:"error_message" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at"`
	ProcessedAt    *time.Time `json:"processed_at"`
	Agent          Agent      `json:"-" gorm:"foreignKey:AgentID"`
}

// AgentStats представляет статистику агента
type AgentStats struct {
	AgentID      uuid.UUID `json:"agent_id"`
	TotalDialogs int64     `json:"total_dialogs"`
	TotalFiles   int64     `json:"total_files"`
	TotalVectors int64     `json:"total_vectors"`
	LastActivity time.Time `json:"last_activity"`
}
