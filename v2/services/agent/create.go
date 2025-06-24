package agent

import (
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/models"
	"macdent-ai-chatbot/v2/utils"
)

type CreateAgentRequest struct {
	APIKey              string             `json:"api_key" validate:"required,min=144"`
	Model               string             `json:"model" validate:"required"`
	SystemPrompt        string             `json:"system_prompt"`
	UserPrompt          string             `json:"user_prompt"`
	ContextSize         int                `json:"context_size"`
	Temperature         float64            `json:"temperature"`
	TopP                float64            `json:"top_p"`
	MaxCompletionTokens int                `json:"max_completion_tokens"`
	Permissions         PermissionsRequest `json:"permissions"`
}

type PermissionsRequest struct {
	Stomatology  bool `json:"stomatology"`
	Doctors      bool `json:"doctors"`
	Appointments bool `json:"appointments"`
}

func (s *Service) CreateAgent(request *CreateAgentRequest, postgres *databases.PostgresDatabase) (*models.Agent, *utils.UserErrorResponse) {
	agent := &models.Agent{
		APIKey:              request.APIKey,
		Model:               request.Model,
		SystemPrompt:        request.SystemPrompt,
		UserPrompt:          request.UserPrompt,
		ContextSize:         request.ContextSize,
		Temperature:         request.Temperature,
		MaxCompletionTokens: request.MaxCompletionTokens,
	}

	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			s.logger.Fatalf("откат транзакций из-за panic: %v", r)
		}
	}()

	if err := tx.Create(agent).Error; err != nil {
		tx.Rollback()

		s.logger.Errorf("создание агента: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка создания агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	permission := &models.Permission{
		AgentID:      agent.ID,
		Stomatology:  request.Permissions.Stomatology,
		Doctors:      request.Permissions.Doctors,
		Appointments: request.Permissions.Appointments,
	}

	if err := tx.Create(permission).Error; err != nil {
		tx.Rollback()

		s.logger.Errorf("создание разрешений: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка создания агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	if err := tx.Model(agent).Association("Permission").Find(&agent.Permission); err != nil {
		tx.Rollback()

		s.logger.Errorf("загрузка разрешений: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка создания агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()

		s.logger.Errorf("закрытие транзакций: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка создания агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	return agent, nil
}
