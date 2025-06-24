package agent

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/models"
	"macdent-ai-chatbot/v2/utils"
)

type GetAgentsRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (s *Service) GetAgent(agentID uuid.UUID, postgres *databases.PostgresDatabase) (*models.Agent, *utils.UserErrorResponse) {
	var agent models.Agent

	err := postgres.DB.
		Preload("Permission").
		Where("id = ?", agentID).
		First(&agent).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Infof("агент с ID %s не найден", agentID)
			return nil, utils.NewUserErrorResponse(
				404,
				"Агент не найден",
				"Указанный агент не существует или был удален.",
			)
		}

		s.logger.Errorf("получение агента: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка получения агента",
			"Пожалуйста, повторите попытку позже",
		)
	}

	return &agent, nil
}

func (s *Service) GetAgents(request *GetAgentsRequest, postgres *databases.PostgresDatabase) ([]*models.Agent, *utils.UserErrorResponse) {
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
		s.logger.Errorf("получение списка агентов: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка получения списка агентов",
			"Пожалуйста, повторите попытку позже",
		)
	}

	return agents, nil
}
