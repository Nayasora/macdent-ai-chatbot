package dialog

import (
	"context"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"macdent-ai-chatbot/internal/databases"
	"macdent-ai-chatbot/internal/models"
	"macdent-ai-chatbot/internal/services/agent"
	openai2 "macdent-ai-chatbot/internal/services/openai"
	"macdent-ai-chatbot/internal/services/tool"
	"macdent-ai-chatbot/internal/utils"
	"time"
)

type UserDialogNewMessageRequest struct {
	AgentID string `json:"agent_id" validate:"required,uuid"`
	UserID  string `json:"user_id" validate:"required"`
	Message string `json:"message" validate:"required,max=1000"`
}

func (s *Service) ResponseDialogNewMessageRequest(
	request *UserDialogNewMessageRequest,
	postgres *databases.PostgresDatabase,
) (string, *utils.UserErrorResponse) {
	agentUUID, _ := uuid.Parse(request.AgentID)

	currentAgent, errorResponse := agent.NewService().
		GetAgent(agentUUID, postgres)

	if errorResponse != nil {
		return "", errorResponse
	}

	var messages []openai.ChatCompletionMessageParamUnion

	if currentAgent.SystemPrompt != "" {
		messages = append(messages, openai.SystemMessage(currentAgent.SystemPrompt))
	}
	if currentAgent.UserPrompt != "" {
		messages = append(messages, openai.SystemMessage(currentAgent.UserPrompt))
	}

	messages = append(messages, openai.UserMessage(request.Message))

	s.logger.Infof("сообщения для OpenAI: %v", messages)

	toolService := tool.NewService(currentAgent)

	return s.processMessagesWithTools(currentAgent, messages, toolService)
}

func (s *Service) processMessagesWithTools(
	agent *models.Agent,
	messages []openai.ChatCompletionMessageParamUnion,
	toolService *tool.Service,
) (string, *utils.UserErrorResponse) {
	chatCompletionParams := s.GetChatCompletionParams(agent, messages)
	chatCompletionParams.Tools = openai.F(toolService.GetToolsFunctions())
	completion, err := s.QueryCompletion(agent, chatCompletionParams)

	if err != nil {
		return "", err
	}

	toolMessage := completion.Choices[0].Message
	s.logger.Infof("ответ OpenAI: %v", toolMessage)

	if toolService.HasToolCalls(toolMessage.ToolCalls) {
		updatedMessages := toolService.ExecuteToolCalls(messages, toolMessage)

		return s.processMessagesWithTools(agent, updatedMessages, toolService)
	}

	return toolMessage.Content, nil
}

func (s *Service) GetChatCompletionParams(agent *models.Agent, messages []openai.ChatCompletionMessageParamUnion) openai.ChatCompletionNewParams {
	return openai.ChatCompletionNewParams{
		Model:               openai.F(agent.Model),
		Messages:            openai.F(messages),
		Temperature:         openai.F(agent.Temperature),
		Modalities:          openai.F([]openai.ChatCompletionModality{openai.ChatCompletionModalityText}),
		MaxCompletionTokens: openai.F(int64(agent.MaxCompletionTokens)),
	}
}

func (s *Service) QueryCompletion(agent *models.Agent, completionParams openai.ChatCompletionNewParams) (*openai.ChatCompletion, *utils.UserErrorResponse) {
	openaiService := openai2.NewService(agent.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	completion, err := openaiService.Client.Chat.Completions.New(ctx, completionParams)

	if err != nil {
		s.logger.Errorf("запрос к OpenAI с результатами инструментов: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка обработки сообщения",
			"Не удалось обработать ваше сообщение. Пожалуйста, попробуйте позже.",
		)
	}

	if len(completion.Choices) == 0 {
		return nil, utils.NewUserErrorResponse(
			500,
			"Пустой ответ",
			"Сервис не смог сформировать ответ на ваше сообщение. Пожалуйста, перефразируйте или попробуйте позже.",
		)
	}

	return completion, nil
}
