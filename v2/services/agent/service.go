package agent

import (
	"github.com/charmbracelet/log"
	"macdent-ai-chatbot/v2/utils"
)

type Service struct {
	logger *log.Logger
}

func NewService() *Service {
	logger := utils.NewLogger("agent")

	return &Service{
		logger: logger,
	}
}
