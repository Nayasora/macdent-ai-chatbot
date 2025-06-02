package store

import (
	"context"
	"crypto/tls"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/config"
	"macdent-ai-chatbot/logger"
	"time"
)

// NewQdrantClient создает новый клиент Qdrant
func NewQdrantClient(cfg *config.QdrantConfig) *qdrant.Client {
	log := logger.New("qdrant")

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   cfg.Host,
		Port:                   cfg.Port,
		APIKey:                 cfg.ApiKey,
		SkipCompatibilityCheck: true,
		UseTLS:                 true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err != nil {
		log.Fatalf("ошибка при создании клиента: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	health, err := client.HealthCheck(ctx)
	if err != nil {
		log.Fatalf("ошибка подключения к базе: %v", err)
	}

	log.Infof("успешное подключение к базе: %s", health)

	return client
}
