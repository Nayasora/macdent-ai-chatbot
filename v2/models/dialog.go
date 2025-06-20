package models

import (
	"github.com/google/uuid"
	"time"
)

type Dialog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AgentID   uuid.UUID `json:"agent_id" gorm:"type:uuid;not null;index"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	Response  string    `json:"response" gorm:"type:text"`
	Role      string    `json:"role" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Agent     Agent     `json:"-" gorm:"foreignKey:AgentID"`
}
