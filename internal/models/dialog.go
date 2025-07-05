package models

import (
	"github.com/google/uuid"
	"time"
)

// Dialog представляет запись диалога между пользователем и агентом
type Dialog struct {
	// Уникальный идентификатор диалога
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Связи
	AgentID uuid.UUID `json:"agent_id" gorm:"type:uuid;not null;index"`
	UserID  string    `json:"user_id" gorm:"not null;index"`

	// Содержание диалога
	Message  string `json:"message" gorm:"type:text;not null"`
	Response string `json:"response" gorm:"type:text"`
	Role     string `json:"role" gorm:"not null;index"`

	// Метаданные
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
}
