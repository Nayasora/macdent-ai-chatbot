package openai

import (
	"github.com/charmbracelet/log"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"macdent-ai-chatbot/v2/utils"
)

type Service struct {
	Client  *openai.Client
	loggger *log.Logger
}

func NewService(apiKey string) *Service {
	customLogger := utils.NewLogger("openai")

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &Service{
		Client:  client,
		loggger: customLogger,
	}
}
