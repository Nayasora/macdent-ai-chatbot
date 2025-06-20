package models

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgentID      uuid.UUID `json:"agent_id" gorm:"type:uuid;uniqueIndex;not null"`
	Stomatology  bool      `json:"stomatology" gorm:"default:false"`
	Doctors      bool      `json:"doctors" gorm:"default:false"`
	Appointments bool      `json:"appointments" gorm:"default:false"`
}
