package tool

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/openai/openai-go"
	"macdent-ai-chatbot/v2/clients"
	"macdent-ai-chatbot/v2/models"
	"macdent-ai-chatbot/v2/utils"
)

type Service struct {
	Agent  *models.Agent
	logger *log.Logger
}

func NewService(agent *models.Agent) *Service {
	logger := utils.NewLogger("tool")

	return &Service{
		Agent:  agent,
		logger: logger,
	}
}

func (s *Service) HasToolCalls(toolCalls []openai.ChatCompletionMessageToolCall) bool {

	if len(toolCalls) == 0 {
		s.logger.Infof("вызов инструментов: не обнаружено")
		return false
	}

	s.logger.Infof("вызов инструментов: обнаружено %d вызовов", len(toolCalls))
	return true
}

func (s *Service) ExecuteToolCalls(messages []openai.ChatCompletionMessageParamUnion, toolMessage openai.ChatCompletionMessage) []openai.ChatCompletionMessageParamUnion {
	var toolResults []openai.ChatCompletionMessageParamUnion

	toolResults = append(toolResults, messages...)
	toolResults = append(toolResults, toolMessage)

	for _, toolCall := range toolMessage.ToolCalls {
		switch toolCall.Function.Name {
		case "get_doctors":
			s.logger.Info("вызов инструмента get_doctors")
			doctorsResponse, errorResponse := clients.GetDoctors(&clients.GetDoctorsRequest{
				AccessToken: s.Agent.Metadata.AccessToken,
			})

			if errorResponse != nil {
				s.logger.Errorf("обработка инструмента get_doctors: %v", errorResponse)
				errorJSON, err := json.Marshal(errorResponse)
				if err != nil {
					s.logger.Errorf("создание json: %v", err)
					continue
				}

				toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(errorJSON)))
				continue
			}

			responseJSON, err := json.Marshal(doctorsResponse)
			if err != nil {
				s.logger.Errorf("создание json: %v", err)
				continue
			}
			toolResults = append(toolResults, openai.ToolMessage(toolCall.ID, string(responseJSON)))
		}
	}

	s.logger.Infof("Выполнение инструментов завершено, добавлено %d сообщений", len(toolResults)-len(messages)-1)
	return toolResults
}
