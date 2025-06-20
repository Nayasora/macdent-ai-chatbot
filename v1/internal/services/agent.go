package services

import (
	"fmt"
	"github.com/google/uuid"
	"macdent-ai-chatbot/v1/internal/database"
	"macdent-ai-chatbot/v1/internal/models"
	"time"
)

type AgentService struct {
	db               *database.Database
	knowledgeService *KnowledgeService
}

type CreateAgentRequest struct {
	ID              *uuid.UUID `json:"id,omitempty"`
	Name            string     `json:"name" validate:"required,min=1,max=100"`
	Description     string     `json:"description,omitempty"`
	APIKey          string     `json:"api_key" validate:"required"`
	Model           string     `json:"model" validate:"required"`
	ContextSize     int        `json:"context_size" validate:"min=1,max=32768"`
	SystemPrompt    string     `json:"system_prompt,omitempty"`
	AssistantPrompt string     `json:"assistant_prompt,omitempty"`
	Temperature     *float32   `json:"temperature,omitempty"`
	MaxTokens       *int       `json:"max_tokens,omitempty"`
}

type UpdateAgentRequest struct {
	Name            *string  `json:"name,omitempty"`
	Description     *string  `json:"description,omitempty"`
	APIKey          *string  `json:"api_key,omitempty"`
	Model           *string  `json:"model,omitempty"`
	ContextSize     *int     `json:"context_size,omitempty"`
	SystemPrompt    *string  `json:"system_prompt,omitempty"`
	AssistantPrompt *string  `json:"assistant_prompt,omitempty"`
	Temperature     *float32 `json:"temperature,omitempty"`
	MaxTokens       *int     `json:"max_tokens,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

func NewAgentService(db *database.Database, knowledgeService *KnowledgeService) *AgentService {
	return &AgentService{
		db:               db,
		knowledgeService: knowledgeService,
	}
}

func (a *AgentService) CreateAgent(req *CreateAgentRequest) (*models.Agent, error) {
	// Генерируем ID если не предоставлен
	agentID := uuid.New()
	if req.ID != nil {
		agentID = *req.ID
	}

	// Проверяем уникальность ID
	existingAgent, _ := a.db.GetAgentByID(agentID.String())
	if existingAgent != nil {
		return nil, fmt.Errorf("agent with ID %s already exists", agentID)
	}

	if req.Model == "" {
		req.Model = "gpt-4.1"
	}

	// Валидируем API ключ
	if err := a.validateOpenAIKey(req.APIKey, req.Model); err != nil {
		return nil, fmt.Errorf("invalid OpenAI API key or model: %w", err)
	}

	// Устанавливаем значения по умолчанию
	temperature := float32(0.7)
	if req.Temperature != nil {
		temperature = *req.Temperature
	}

	maxTokens := 1000
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	contextSize := 4096
	if req.ContextSize > 0 {
		contextSize = req.ContextSize
	}

	agent := &models.Agent{
		ID:              agentID,
		Name:            req.Name,
		Description:     req.Description,
		APIKey:          req.APIKey,
		Model:           req.Model,
		ContextSize:     contextSize,
		SystemPrompt:    req.SystemPrompt,
		AssistantPrompt: req.AssistantPrompt,
		Temperature:     temperature,
		MaxTokens:       maxTokens,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Создаем агента в базе данных
	err := a.db.CreateAgent(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Инициализируем базу знаний для агента
	err = a.knowledgeService.InitializeAgentKnowledge(agentID)
	if err != nil {
		// Откатываем создание агента
		a.db.DeleteAgent(agentID.String())
		return nil, fmt.Errorf("failed to initialize knowledge base: %w", err)
	}

	return agent, nil
}

func (a *AgentService) GetAgent(agentID uuid.UUID) (*models.Agent, error) {
	agent, err := a.db.GetAgentByID(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	if !agent.IsActive {
		return nil, fmt.Errorf("agent is deactivated")
	}

	return agent, nil
}

func (a *AgentService) UpdateAgent(agentID uuid.UUID, req *UpdateAgentRequest) (*models.Agent, error) {
	agent, err := a.db.GetAgentByID(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// Валидируем новый API ключ если предоставлен
	if req.APIKey != nil {
		model := agent.Model
		if req.Model != nil {
			model = *req.Model
		}
		if err := a.validateOpenAIKey(*req.APIKey, model); err != nil {
			return nil, fmt.Errorf("invalid OpenAI API key or model: %w", err)
		}
	}

	// Обновляем только предоставленные поля
	if req.Name != nil {
		agent.Name = *req.Name
	}
	if req.Description != nil {
		agent.Description = *req.Description
	}
	if req.APIKey != nil {
		agent.APIKey = *req.APIKey
	}
	if req.Model != nil {
		agent.Model = *req.Model
	}
	if req.ContextSize != nil {
		agent.ContextSize = *req.ContextSize
	}
	if req.SystemPrompt != nil {
		agent.SystemPrompt = *req.SystemPrompt
	}
	if req.AssistantPrompt != nil {
		agent.AssistantPrompt = *req.AssistantPrompt
	}
	if req.Temperature != nil {
		agent.Temperature = *req.Temperature
	}
	if req.MaxTokens != nil {
		agent.MaxTokens = *req.MaxTokens
	}
	if req.IsActive != nil {
		agent.IsActive = *req.IsActive
	}

	agent.UpdatedAt = time.Now()

	err = a.db.UpdateAgent(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}

	return agent, nil
}

func (a *AgentService) DeleteAgent(agentID uuid.UUID) error {
	// Проверяем существование агента
	_, err := a.db.GetAgentByID(agentID.String())
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Удаляем все знания агента
	err = a.knowledgeService.DeleteAllAgentKnowledge(agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent knowledge: %w", err)
	}

	// Удаляем агента
	err = a.db.DeleteAgent(agentID.String())
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	return nil
}

func (a *AgentService) ListAgents(limit, offset int) ([]models.Agent, error) {
	return a.db.ListAgents(limit, offset)
}

func (a *AgentService) GetAgentStats(agentID uuid.UUID) (*models.AgentStats, error) {
	agent, err := a.db.GetAgentByID(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// Получаем статистику диалогов
	totalDialogs, err := a.db.CountDialogsByAgent(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to count dialogs: %w", err)
	}

	// Получаем статистику файлов
	files, err := a.knowledgeService.GetKnowledgeFiles(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge files: %w", err)
	}

	// Получаем статистику векторов
	knowledgeStats, err := a.knowledgeService.GetKnowledgeStats(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge stats: %w", err)
	}

	// Получаем последнюю активность
	lastActivity, err := a.db.GetLastDialogTime(agentID.String())
	if err != nil {
		lastActivity = agent.CreatedAt
	}

	vectorCount := int64(0)
	if count, ok := knowledgeStats["vector_count"].(uint64); ok {
		vectorCount = int64(count)
	}

	return &models.AgentStats{
		AgentID:      agentID,
		TotalDialogs: totalDialogs,
		TotalFiles:   int64(len(files)),
		TotalVectors: vectorCount,
		LastActivity: lastActivity,
	}, nil
}

func (a *AgentService) validateOpenAIKey(apiKey, model string) error {
	// Создаем временный OpenAI сервис для проверки ключа
	openaiService, err := NewOpenAIService(apiKey)
	if err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}

	// Проверяем ключ
	err = openaiService.ValidateAPIKey()
	if err != nil {
		return fmt.Errorf("API key validation failed: %w", err)
	}

	// Проверяем поддержку модели
	supportedModels := openaiService.GetSupportedModels()
	modelSupported := false
	for _, supportedModel := range supportedModels {
		if supportedModel == model {
			modelSupported = true
			break
		}
	}

	if !modelSupported {
		return fmt.Errorf("model %s is not supported", model)
	}

	return nil
}

func (a *AgentService) GetSupportedModels() []string {
	// Создаем временный сервис для получения списка моделей
	tempService, _ := NewOpenAIService("dummy-key")
	if tempService != nil {
		return tempService.GetSupportedModels()
	}

	// Fallback список
	return []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	}
}
