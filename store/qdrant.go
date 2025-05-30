package store

import (
	"context"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/config"
	"macdent-ai-chatbot/logger"
)

// NewQdrantClient создает новый клиент Qdrant
func NewQdrantClient(cfg *config.QdrantConfig) *qdrant.Client {
	log := logger.New("qdrant")

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   cfg.Host,
		Port:                   cfg.Port,
		SkipCompatibilityCheck: true,
	})

	if err != nil {
		log.Fatalf("ошибка при создании клиента: %v", err)
	}

	health, err := client.HealthCheck(context.Background())
	if err != nil {
		log.Fatalf("ошибка подключения к базе: %v", err)
	}

	log.Infof("успешное подключение к базе: %s", health)

	return client
}
