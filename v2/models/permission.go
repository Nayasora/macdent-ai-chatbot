package models

import (
	"github.com/google/uuid"
	"time"
)

// Permission определяет права доступа агента к различным модулям
type Permission struct {
	// Идентификатор агента, к которому относятся эти разрешения
	AgentID uuid.UUID `json:"agent_id" gorm:"type:uuid;primary_key;not null"`

	// Разрешения по модулям
	Stomatology  bool `json:"stomatology" gorm:"default:false;not null"`
	Doctors      bool `json:"doctors" gorm:"default:false;not null"`
	Appointments bool `json:"appointments" gorm:"default:false;not null"`

	// Метаданные
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
}
