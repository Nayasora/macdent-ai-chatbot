package openai

import (
	"context"
	"strings"
	"time"
)

type GetModelsRequest struct {
	ApiKey string `json:"api_key" validate:"required,min=144"`
}

type GetModelsErrorResponse struct {
	StatusCode int
	Message    string
	Details    string
}

func (s *Service) GetModels() ([]string, *GetModelsErrorResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := s.client.Models.List(ctx)

	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "invalid_api_key") {
			return nil, &GetModelsErrorResponse{
				StatusCode: 400,
				Message:    "Неверный API ключ",
			}
		}

		s.log.Fatalf("получение списка моделей %s", err)
	}

	var names []string

	for _, model := range list.Data {
		if strings.Contains(strings.ToLower(model.ID), "gpt") {
			names = append(names, model.ID)
		}
	}

	return names, nil
}
