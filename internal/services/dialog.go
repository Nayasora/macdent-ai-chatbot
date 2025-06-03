package services

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"macdent-ai-chatbot/internal/database"
	"macdent-ai-chatbot/internal/models"
	"time"
)

type DialogService struct {
	db               *database.Database
	knowledgeService *KnowledgeService
}

type DialogRequest struct {
	Message        string                 `json:"message" validate:"required,min=1"`
	UserID         string                 `json:"user_id" validate:"required"`
	UseKnowledge   bool                   `json:"use_knowledge" default:"true"`
	KnowledgeLimit int                    `json:"knowledge_limit" default:"3"`
	ScoreThreshold float32                `json:"score_threshold" default:"0.3"`
	HistoryLimit   int                    `json:"history_limit" default:"10"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type DialogResponse struct {
	Response      string                  `json:"response"`
	AgentID       string                  `json:"agent_id"`
	DialogID      uint                    `json:"dialog_id"`
	KnowledgeUsed []KnowledgeSearchResult `json:"knowledge_used,omitempty"`
	ResponseTime  time.Duration           `json:"response_time"`
	TokensUsed    int                     `json:"tokens_used,omitempty"`
	Model         string                  `json:"model,omitempty"`
	Metadata      map[string]interface{}  `json:"metadata,omitempty"`
}

func NewDialogService(db *database.Database, knowledgeService *KnowledgeService) *DialogService {
	return &DialogService{
		db:               db,
		knowledgeService: knowledgeService,
	}
}

func (d *DialogService) ProcessDialog(agentID uuid.UUID, req *DialogRequest) (*DialogResponse, error) {
	startTime := time.Now()

	// Получаем агента
	agent, err := d.db.GetAgentByID(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	if !agent.IsActive {
		return nil, fmt.Errorf("agent is deactivated")
	}

	// Устанавливаем значения по умолчанию
	if req.KnowledgeLimit <= 0 {
		req.KnowledgeLimit = 3
	}
	if req.ScoreThreshold <= 0 {
		req.ScoreThreshold = 0.5
	}
	if req.HistoryLimit <= 0 {
		req.HistoryLimit = 10
	}

	// Сохраняем сообщение пользователя
	_, err = d.saveUserMessage(agentID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// Получаем историю диалогов
	dialogHistory, err := d.db.GetDialogHistory(agentID.String(), req.UserID, req.HistoryLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get dialog history: %w", err)
	}

	// Поиск в базе знаний (RAG)
	var knowledgeContext []KnowledgeSearchResult
	if req.UseKnowledge {
		knowledgeContext, err = d.knowledgeService.SearchKnowledge(
			agentID,
			req.Message,
			req.KnowledgeLimit,
			req.ScoreThreshold,
		)
		if err != nil {
			fmt.Printf("Knowledge search failed: %v\n", err)
			knowledgeContext = []KnowledgeSearchResult{}
		}
	}

	// Создаем OpenAI сервис с API ключом агента
	openaiService, err := NewOpenAIService(agent.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI service: %w", err)
	}

	fmt.Println("База знаний:", knowledgeContext)
	// Генерируем ответ
	chatResponse, err := openaiService.GenerateResponse(*agent, req.Message, dialogHistory, knowledgeContext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Сохраняем ответ ассистента
	assistantDialog, err := d.saveAssistantResponse(agentID, req.UserID, chatResponse.Message, req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to save assistant response: %w", err)
	}

	responseTime := time.Since(startTime)

	return &DialogResponse{
		Response:      chatResponse.Message,
		AgentID:       agentID.String(),
		DialogID:      assistantDialog.ID,
		KnowledgeUsed: knowledgeContext,
		ResponseTime:  responseTime,
		TokensUsed:    chatResponse.TokensUsed,
		Model:         chatResponse.Model,
		Metadata:      req.Metadata,
	}, nil
}

func (d *DialogService) GetDialogHistory(agentID uuid.UUID, userID string, limit int) ([]models.Dialog, error) {
	if limit <= 0 {
		limit = 50
	}

	return d.db.GetDialogHistory(agentID.String(), userID, limit)
}

func (d *DialogService) DeleteDialogHistory(agentID uuid.UUID, userID string) error {
	return d.db.DeleteDialogHistory(agentID.String(), userID)
}

func (d *DialogService) GetDialogsByAgent(agentID uuid.UUID, limit, offset int) ([]models.Dialog, error) {
	return d.db.GetDialogsByAgent(agentID.String(), limit, offset)
}

func (d *DialogService) GetDialogStats(agentID uuid.UUID) (map[string]interface{}, error) {
	totalDialogs, err := d.db.CountDialogsByAgent(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to count dialogs: %w", err)
	}

	lastActivity, err := d.db.GetLastDialogTime(agentID.String())
	if err != nil {
		lastActivity = time.Time{}
	}

	return map[string]interface{}{
		"total_dialogs": totalDialogs,
		"last_activity": lastActivity,
		"agent_id":      agentID.String(),
	}, nil
}

func (d *DialogService) saveUserMessage(agentID uuid.UUID, req *DialogRequest) (*models.Dialog, error) {
	metadataJSON := "{}"
	if req.Metadata != nil {
		metadataBytes, _ := json.Marshal(req.Metadata)
		metadataJSON = string(metadataBytes)
	}

	userDialog := &models.Dialog{
		AgentID:   agentID,
		UserID:    req.UserID,
		Message:   req.Message,
		Role:      "user",
		Metadata:  metadataJSON,
		CreatedAt: time.Now(),
	}

	err := d.db.CreateDialog(userDialog)
	if err != nil {
		return nil, err
	}

	return userDialog, nil
}

func (d *DialogService) saveAssistantResponse(agentID uuid.UUID, userID, response string, metadata map[string]interface{}) (*models.Dialog, error) {
	metadataJSON := "{}"
	if metadata != nil {
		metadataBytes, _ := json.Marshal(metadata)
		metadataJSON = string(metadataBytes)
	}

	assistantDialog := &models.Dialog{
		AgentID:   agentID,
		UserID:    userID,
		Message:   response,
		Role:      "assistant",
		Metadata:  metadataJSON,
		CreatedAt: time.Now(),
	}

	err := d.db.CreateDialog(assistantDialog)
	if err != nil {
		return nil, err
	}

	return assistantDialog, nil
}
