package databases

import (
	"crypto/tls"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/v2/configs"
	"macdent-ai-chatbot/v2/utils"
)

type QdrantDatabase struct {
	client *qdrant.Client
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
		client: client,
	}
}
