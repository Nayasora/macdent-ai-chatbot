package models

import (
	"github.com/google/uuid"
	"time"
)

type KnowledgePrompt struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AgentID   uuid.UUID `json:"agent_id" gorm:"type:uuid;not null;index"`
	Prompt    string    `json:"prompt" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	CreatedAt      time.Time  `json:"created_at"`
	ProcessedAt    *time.Time `json:"processed_at"`
	Agent          Agent      `json:"-" gorm:"foreignKey:AgentID"`
}
