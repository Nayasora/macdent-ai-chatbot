package agent

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/models"
	"macdent-ai-chatbot/v2/utils"
)

type UpdateAgentRequest struct {
	AgentID             string             `json:"agent_id" validate:"required,uuid"`
	APIKey              string             `json:"api_key"`
	Model               string             `json:"model"`
	SystemPrompt        string             `json:"system_prompt"`
	UserPrompt          string             `json:"user_prompt"`
	ContextSize         *int               `json:"context_size"`
	Temperature         *float64           `json:"temperature"`
	TopP                *float64           `json:"top_p"`
	MaxCompletionTokens *int               `json:"max_completion_tokens"`
	Metadata            *MetadataRequest   `json:"metadata"`
	Permissions         PermissionsRequest `json:"permissions"`
}

func (s *Service) UpdateAgent(request *UpdateAgentRequest, postgres *databases.PostgresDatabase) (*models.Agent, *utils.UserErrorResponse) {
	agentID, _ := uuid.Parse(request.AgentID)
	agent, errorResponse := s.GetAgent(
		agentID,
		postgres,
	)

	if errorResponse != nil {
		return nil, errorResponse
	}

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

	if request.ContextSize != nil {
		agent.ContextSize = *request.ContextSize
	}
	if request.Temperature != nil {
		agent.Temperature = *request.Temperature
	}
	if request.MaxCompletionTokens != nil {
		agent.MaxCompletionTokens = *request.MaxCompletionTokens
	}

	if request.Metadata != nil {
		if request.Metadata.Stomatology != 0 {
			agent.Metadata.Stomatology = request.Metadata.Stomatology
		}
		if request.Metadata.AccessToken != "" {
			agent.Metadata.AccessToken = request.Metadata.AccessToken
		}
	}

	agent.Permission.Stomatology = request.Permissions.Stomatology
	agent.Permission.Doctors = request.Permissions.Doctors
	agent.Permission.Appointment = request.Permissions.Appointment

	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			s.logger.Fatalf("откат транзакций из-за panic: %v", r)
		}
	}()

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(agent).Error; err != nil {
		tx.Rollback()

		s.logger.Errorf("обновление агента: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка обновления агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()

		s.logger.Errorf("закрытие транзакций: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка обновления агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	return agent, nil
}
