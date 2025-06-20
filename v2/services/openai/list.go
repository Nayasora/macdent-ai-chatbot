package openai

import (
	"context"
	"strings"
	"time"
)

type GetModelsRequest struct {
	ApiKey string `json:"api_key" validate:"required,min=144"`
}

func (s *Service) GetModels() []string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := s.client.Models.List(ctx)

	if err != nil {
		s.log.Fatalf("получение списка моделей %s", err)
	}

	var names []string

	for _, model := range list.Data {
		if strings.Contains(strings.ToLower(model.ID), "gpt") {
			names = append(names, model.ID)
		}
	}

	return names
}
