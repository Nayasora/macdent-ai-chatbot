package agent

import (
	"github.com/google/uuid"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/models"
)

type GetAgentRequest struct {
	AgentID string `json:"agent_id"`
}

type GetAgentsRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (s *Service) GetAgent(request *GetAgentRequest, postgres *databases.PostgresDatabase) *models.Agent {
	var agent models.Agent

	// Преобразуем строковый ID в UUID
	id, err := uuid.Parse(request.AgentID)
	if err != nil {
		s.logger.Fatalf("неверный формат ID агента: %v", err)
	}

	err = postgres.DB.
		Preload("Permission").
		Where("id = ?", id).
		First(&agent).Error

	if err != nil {
		s.logger.Fatalf("получение агента: %v", err)
	}

	return &agent
}

func (s *Service) GetAgents(request *GetAgentsRequest, postgres *databases.PostgresDatabase) []*models.Agent {
	var agents []*models.Agent

	if request.Limit <= 0 {
		request.Limit = 10
	}

	err := postgres.DB.
		Preload("Permission").
		Order("created_at DESC").
		Limit(request.Limit).
		Offset(request.Offset).
		Find(&agents).Error

	if err != nil {
		s.logger.Fatalf("получение списка агентов: %v", err)
	}

	return agents
}
