package knowledge

import (
	"context"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/internal/utils"
)

func (s *Service) CreateCollection(ctx context.Context, options *qdrant.CreateCollection) *utils.UserErrorResponse {
	collectionExists, err := s.qdrant.Client.CollectionExists(ctx, options.CollectionName)
	if err != nil {
		s.logger.Errorf("проверка существования коллекции в Qdrant: %v", err)
		return utils.NewUserErrorResponse(
			500,
			"Ошибка загрузки базы знаний",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	if !collectionExists {
		err = s.qdrant.Client.CreateCollection(ctx, options)
		if err != nil {
			s.logger.Errorf("создание коллекции в Qdrant: %v", err)
			return utils.NewUserErrorResponse(
				500,
				"Ошибка загрузки базы знаний",
				"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
			)
		}
		s.logger.Infof("создана коллекция %s в Qdrant", options.CollectionName)
	}

	return nil
}
