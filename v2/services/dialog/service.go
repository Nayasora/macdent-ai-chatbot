package dialog

import (
	"github.com/charmbracelet/log"
	"macdent-ai-chatbot/v2/utils"
)

type Service struct {
	logger *log.Logger
}

func NewService() *Service {
	logger := utils.NewLogger("dialog")

	return &Service{
		logger: logger,
	}
}
