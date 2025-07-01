package openai

import (
	"context"
	"macdent-ai-chatbot/v2/utils"
	"strings"
	"time"
)

type GetModelsRequest struct {
	ApiKey string `json:"api_key" validate:"required,min=144"`
}

func (s *Service) GetModels() ([]string, *utils.UserErrorResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := s.Client.Models.List(ctx)

	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "invalid_api_key") {
			return nil, utils.NewUserErrorResponse(
				400,
				"Неверный API ключ",
				"Пожалуйста, проверьте ваш API ключ и попробуйте снова.",
			)
		}

		s.logger.Errorf("получение списка моделей %s", err)

		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка получения моделей",
			"Произошла ошибка при получении списка моделей. Пожалуйста, попробуйте позже.",
		)
	}

	var names []string

	for _, model := range list.Data {
		if strings.Contains(strings.ToLower(model.ID), "gpt") {
			names = append(names, model.ID)
		}
	}

	return names, nil
}
