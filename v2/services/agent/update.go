package agent

import (
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/models"
)

type UpdateAgentRequest struct {
	AgentID      string             `json:"agent_id" validate:"required"`
	APIKey       string             `json:"api_key" gorm:"not null"`
	Model        string             `json:"model" gorm:"not null"`
	SystemPrompt string             `json:"system_prompt" gorm:"not null"`
	UserPrompt   string             `json:"user_prompt" gorm:"not null"`
	ContextSize  int                `json:"context_size" gorm:"not null"`
	Temperature  float32            `json:"temperature" gorm:"not null"`
	TopP         float32            `json:"top_p" gorm:"not null"`
	MaxTokens    int                `json:"max_tokens" gorm:"not null"`
	Permissions  PermissionsRequest `json:"permissions"`
}

func (s *Service) UpdateAgent(request *UpdateAgentRequest, postgres *databases.PostgresDatabase) *models.Agent {
	agent := s.GetAgent(
		&GetAgentRequest{
			AgentID: request.AgentID,
		},
		postgres,
	)

	if request.Model != "" {
		agent.Model = request.Model
	}
	if request.APIKey != "" {
		agent.APIKey = request.APIKey
	}
	if request.SystemPrompt != "" {
		agent.SystemPrompt = request.SystemPrompt
	}
	if request.UserPrompt != "" {
		agent.UserPrompt = request.UserPrompt
	}
	if request.ContextSize > 0 {
		agent.ContextSize = request.ContextSize
	}
	if request.Temperature >= 0.0 {
		agent.Temperature = request.Temperature
	}
	if request.TopP >= 0.0 {
		agent.TopP = request.TopP
	}
	if request.MaxTokens > 0 {
		agent.MaxTokens = request.MaxTokens
	}

	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			s.logger.Fatalf("откат транзакций из-за panic: %v", r)
		}
	}()

	if agent.Permissions == nil {
		agent.Permissions = &models.Permission{
			AgentID:      agent.ID,
			Stomatology:  request.Permissions.Stomatology,
			Doctors:      request.Permissions.Doctors,
			Appointments: request.Permissions.Appointments,
		}
	} else {
		agent.Permissions.Stomatology = request.Permissions.Stomatology
		agent.Permissions.Doctors = request.Permissions.Doctors
		agent.Permissions.Appointments = request.Permissions.Appointments
	}

	if err := tx.Save(agent).Error; err != nil {
		tx.Rollback()
		s.logger.Fatalf("обновление агента: %v", err)
		return nil
	}

	if err := tx.Save(agent.Permissions).Error; err != nil {
		tx.Rollback()
		s.logger.Fatalf("обновление разрешений агента: %v", err)
		return nil
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		s.logger.Fatalf("закрытие транзакций: %v", err)
		return nil
	}

	return agent
}
