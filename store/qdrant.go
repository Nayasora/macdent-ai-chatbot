package store

import (
	"context"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/logger"
	"os"
)

type QdrantConfig struct {
	Host string
	Port int
}

func New(config *QdrantConfig) *qdrant.Client {
	log := logger.New("qdrant")

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   config.Host,
		Port:                   config.Port,
		SkipCompatibilityCheck: true,
	})

	if err != nil {
		log.Errorf("ошибка при создании клиента: %v", err)
		os.Exit(1)
	}

	health, err := client.HealthCheck(context.Background())
	if err != nil {
		log.Errorf("ошибка подключения к базе: %v", err)
		os.Exit(1)
	}

	log.Infof("успешное подключение к базе: %s", health)

	return client
}
