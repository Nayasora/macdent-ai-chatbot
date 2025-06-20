package agent

import (
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
	err := postgres.DB.
		Preload("Permission").
		Where("id = ?", request.AgentID).
		First(&agent).Error

	if err != nil {
		s.logger.Fatalf("получение агента: %v", err)
	}

	return &agent
}

func (s *Service) GetAgents(request *GetAgentsRequest, postgres *databases.PostgresDatabase) []*models.Agent {
	var agents []*models.Agent
	err := postgres.DB.Order("created_at DESC").
		Limit(request.Limit).
		Offset(request.Offset).
		Find(&agents).Error
	if err != nil {
		s.logger.Fatalf("получение списка агентов: %v", err)
	}
	return agents
}
