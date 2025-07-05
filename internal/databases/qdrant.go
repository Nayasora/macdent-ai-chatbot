package databases

import (
	"crypto/tls"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/internal/configs"
	"macdent-ai-chatbot/internal/utils"
)

type QdrantDatabase struct {
	Client *qdrant.Client
}

func NewQdrant(config *configs.QdrantConfig) *QdrantDatabase {
	logger := utils.NewLogger("qdrant")

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   config.Host,
		Port:                   config.Port,
		UseTLS:                 true,
		SkipCompatibilityCheck: true,
		APIKey:                 config.ApiKey,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	if err != nil {
		logger.Errorf("qdrant клиент: %v", err)
	}

	return &QdrantDatabase{
		Client: client,
	}
}
