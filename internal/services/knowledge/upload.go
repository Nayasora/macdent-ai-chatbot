package knowledge

import (
	"context"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	"macdent-ai-chatbot/internal/models"
	"macdent-ai-chatbot/internal/services/agent"
	openai2 "macdent-ai-chatbot/internal/services/openai"
	"macdent-ai-chatbot/internal/utils"
	"time"
)

type UploadKnowledgeRequest struct {
	AgentID string      `json:"agent_id" validate:"required,uuid"`
	Files   []Knowledge `json:"files" validate:"required,dive,required"`
}

type Knowledge struct {
	Name    string
	Size    int64
	Type    string
	Content []byte
}

func (s *Service) UploadKnowledge(request *UploadKnowledgeRequest) *utils.UserErrorResponse {
	agentUUID, _ := uuid.Parse(request.AgentID)
	currentAgent, errorResponse := agent.NewService().
		GetAgent(agentUUID, s.postgres)

	if errorResponse != nil {
		s.logger.Errorf("получения агента %s: %v", request.AgentID, errorResponse)
		return errorResponse
	}

	// Извлечение информации из файлов
	var knowledgeSize int
	var knowledgeContent []byte
	for _, knowledge := range request.Files {
		knowledgeSize += int(knowledge.Size)
		knowledgeContent = append(knowledgeContent, knowledge.Content...)
	}

	// Создание prompt knowledge
	s.logger.Infof("длина knowledge %d байт", knowledgeSize)
	if knowledgeSize < PromptTypeSize {
		s.CreatePromptKnowledge(string(knowledgeContent), agentUUID)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Создание коллекции в Qdrant
	errorResponse = s.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: request.AgentID,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     3072,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if errorResponse != nil {
		return errorResponse
	}

	openaiService := openai2.NewService(currentAgent.APIKey)
	results, errorResponse := s.PrepareContent(openaiService, string(knowledgeContent), agentUUID)
	if errorResponse != nil {
		return errorResponse
	}

	errorResponse = s.UpsertChunks(ctx, agentUUID, results)
	if errorResponse != nil {
		return errorResponse
	}

	s.CreateKnowledgeFile(agentUUID, request.Files, len(results))

	s.logger.Infof("успешно загружено %d чанков для агента %s", len(results), request.AgentID)
	return nil
}

func (s *Service) UpsertPoints(ctx context.Context, agentID uuid.UUID, points []*qdrant.PointStruct) *utils.UserErrorResponse {
	_, err := s.qdrant.Client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: agentID.String(),
		Points:         points,
	})
	if err != nil {
		s.logger.Errorf("загрузка поинтов в Qdrant: %v", err)
		return utils.NewUserErrorResponse(
			500,
			"Ошибка загрузки базы знаний",
			"Пожалуйста, попробуйте позже или обратитесь в службу поддержки.",
		)
	}

	s.logger.Infof("успешная загрузка поинтов в Qdrant для агента %s", agentID.String())
	return nil
}

func (s *Service) CreatePromptKnowledge(content string, agentUUID uuid.UUID) {
	knowledgePrompt := models.KnowledgePrompt{
		AgentID: agentUUID,
		Prompt:  content,
	}
	s.postgres.DB.Create(&knowledgePrompt)
	s.logger.Infof("создание prompt knowledge для агента %s", agentUUID.String())
	s.logger.Infof("контент knowledge: %s", knowledgePrompt.Prompt)
}

func (s *Service) CreateKnowledgeFile(agentID uuid.UUID, files []Knowledge, chunkCount int) {
	for _, file := range files {
		knowledgeFile := models.KnowledgeFile{
			AgentID:        agentID,
			FileName:       file.Name,
			OriginalName:   file.Name,
			FileSize:       file.Size,
			FileType:       file.Type,
			CollectionName: agentID.String(),
			ChunkCount:     chunkCount,
			Status:         "completed",
			ProcessedAt:    &time.Time{},
		}
		*knowledgeFile.ProcessedAt = time.Now()
		s.postgres.DB.Create(&knowledgeFile)
	}
}
