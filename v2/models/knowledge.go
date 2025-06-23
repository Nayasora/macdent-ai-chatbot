package models

import (
	"github.com/google/uuid"
	"time"
)

// KnowledgePrompt представляет подсказку, используемую в базе знаний агента
type KnowledgePrompt struct {
	// Уникальный идентификатор промпта
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Связь с агентом
	AgentID uuid.UUID `json:"agent_id" gorm:"type:uuid;not null;index"`

	// Содержание промпта
	Prompt string `json:"prompt" gorm:"type:text;not null"`

	// Метаданные
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
}

// KnowledgeFile представляет загруженный и обработанный файл базы знаний
type KnowledgeFile struct {
	// Уникальный идентификатор файла
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Связь с агентом
	AgentID uuid.UUID `json:"agent_id" gorm:"type:uuid;not null;index"`

	// Информация о файле
	FileName     string `json:"file_name" gorm:"not null"`
	OriginalName string `json:"original_name" gorm:"not null"`
	FileSize     int64  `json:"file_size" gorm:"not null"`
	FileType     string `json:"file_type" gorm:"not null"`
	FilePath     string `json:"file_path" gorm:"not null"`

	// Параметры обработки и статус
	CollectionName string `json:"collection_name" gorm:"not null;index"`
	ChunkCount     int    `json:"chunk_count" gorm:"default:0;not null"`
	Status         string `json:"status" gorm:"default:processing;not null;index"`

	// Метаданные
	CreatedAt   time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"not null;autoUpdateTime"`
	ProcessedAt *time.Time `json:"processed_at"`
}
