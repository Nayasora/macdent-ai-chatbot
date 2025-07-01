package knowledge

import (
	"github.com/charmbracelet/log"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/utils"
)

const PromptTypeSize = 3000

type Service struct {
	logger   *log.Logger
	postgres *databases.PostgresDatabase
	qdrant   *databases.QdrantDatabase
}

func NewService(postgres *databases.PostgresDatabase, qdrant *databases.QdrantDatabase) *Service {
	logger := utils.NewLogger("knowledge")

	return &Service{
		logger:   logger,
		postgres: postgres,
		qdrant:   qdrant,
	}
}
