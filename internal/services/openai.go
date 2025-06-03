package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"macdent-ai-chatbot/internal/models"
)

type OpenAIService struct {
	client *openai.Client
	apiKey string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Message      string        `json:"message"`
	Model        string        `json:"model"`
	TokensUsed   int           `json:"tokens_used"`
	ResponseTime time.Duration `json:"response_time"`
}

func NewOpenAIService(apiKey string) (*OpenAIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAIService{
		client: client,
		apiKey: apiKey,
	}, nil
}

func (o *OpenAIService) GenerateResponse(agent models.Agent, userMessage string, dialogHistory []models.Dialog, knowledgeContext []KnowledgeSearchResult) (*ChatResponse, error) {
	startTime := time.Now()

	// Валидируем агента
	if agent.APIKey != o.apiKey {
		return nil, fmt.Errorf("API key mismatch")
	}

	// Строим сообщения для чата
	messages := o.buildChatMessages(agent, userMessage, dialogHistory, knowledgeContext)

	// Создаем запрос к OpenAI
	resp, err := o.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model:               openai.F(agent.Model),
		Messages:            openai.F(messages),
		MaxCompletionTokens: openai.F(int64(agent.MaxTokens)),
		Temperature:         openai.F(float64(agent.Temperature)),
		TopP:                openai.F(1.0),
		FrequencyPenalty:    openai.F(0.0),
		PresencePenalty:     openai.F(0.0),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	responseTime := time.Since(startTime)

	return &ChatResponse{
		Message:      resp.Choices[0].Message.Content,
		Model:        resp.Model,
		TokensUsed:   int(resp.Usage.TotalTokens),
		ResponseTime: responseTime,
	}, nil
}

func (o *OpenAIService) ValidateAPIKey() error {
	// Проверяем API ключ простым запросом
	_, err := o.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.F("gpt-3.5-turbo"),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("test"),
		}),
		MaxCompletionTokens: openai.F(int64(1)),
	})

	return err
}

func (o *OpenAIService) buildChatMessages(agent models.Agent, userMessage string, dialogHistory []models.Dialog, knowledgeContext []KnowledgeSearchResult) []openai.ChatCompletionMessageParamUnion {
	var messages []openai.ChatCompletionMessageParamUnion

	// Добавляем системный промпт
	systemContent := o.buildSystemPrompt(agent, knowledgeContext)
	if systemContent != "" {
		messages = append(messages, openai.SystemMessage(systemContent))
	}

	// Добавляем историю диалогов
	for _, dialog := range dialogHistory {
		switch dialog.Role {
		case "user":
			messages = append(messages, openai.UserMessage(dialog.Message))
		case "assistant":
			messages = append(messages, openai.AssistantMessage(dialog.Message))
		}
	}

	// Добавляем текущее сообщение пользователя
	messages = append(messages, openai.UserMessage(userMessage))

	return messages
}

func (o *OpenAIService) buildSystemPrompt(agent models.Agent, knowledgeContext []KnowledgeSearchResult) string {
	var promptBuilder strings.Builder

	// Добавляем основной системный промпт
	if agent.SystemPrompt != "" {
		promptBuilder.WriteString(agent.SystemPrompt)
		promptBuilder.WriteString("\n\n")
	}

	// Добавляем контекст из базы знаний
	if len(knowledgeContext) > 0 {
		promptBuilder.WriteString("You have access to the following relevant information from the knowledge base:\n\n")

		for i, knowledge := range knowledgeContext {
			promptBuilder.WriteString(fmt.Sprintf("Knowledge %d (Source: %s, Relevance: %.2f):\n",
				i+1, knowledge.Source, knowledge.Score))
			promptBuilder.WriteString(knowledge.Content)
			promptBuilder.WriteString("\n\n")
		}

		promptBuilder.WriteString("Please use this information to provide accurate and helpful responses. If the knowledge base doesn't contain relevant information, you can still provide general assistance.\n\n")
	}

	// Добавляем инструкции для ассистента
	if agent.AssistantPrompt != "" {
		promptBuilder.WriteString("Additional Instructions:\n")
		promptBuilder.WriteString(agent.AssistantPrompt)
		promptBuilder.WriteString("\n\n")
	}

	return strings.TrimSpace(promptBuilder.String())
}

func (o *OpenAIService) GetSupportedModels() []string {
	modelsList, err := o.client.Models.List(context.Background())
	if err != nil {
		fmt.Println("Error fetching models:", err)
		os.Exit(1)
	}

	var modelNames []string

	for _, model := range modelsList.Data {
		modelNames = append(modelNames, model.ID)
	}

	return modelNames
}
