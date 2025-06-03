package services

import (
	"fmt"
	"macdent-ai-chatbot/internal/database"
	"macdent-ai-chatbot/internal/models"
	"time"

	"github.com/google/uuid"
)

type KnowledgeService struct {
	db       *database.Database
	vectorDB *database.VectorDatabase
}

type KnowledgeSearchResult struct {
	Content  string
	Source   string
	Score    float32
	Metadata map[string]interface{}
}

func NewKnowledgeService(db *database.Database, vectorDB *database.VectorDatabase) *KnowledgeService {
	return &KnowledgeService{
		db:       db,
		vectorDB: vectorDB,
	}
}

func (k *KnowledgeService) InitializeAgentKnowledge(agentID uuid.UUID) error {
	collectionName := k.getCollectionName(agentID)

	// Проверяем существование коллекции
	exists, err := k.vectorDB.CollectionExists(collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		// Используем стандартный размер для text-embedding-3-small
		err = k.vectorDB.CreateCollection(collectionName, 1536)
		if err != nil {
			return fmt.Errorf("failed to create collection for agent %s: %w", agentID, err)
		}
	}

	return nil
}

func (k *KnowledgeService) AddDocument(agentID uuid.UUID, fileName string, content string, fileSize int64, fileType string) error {
	// Получаем агента для использования его API ключа
	agent, err := k.db.GetAgentByID(agentID.String())
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Создаем сервис эмбеддингов с API ключом агента
	embeddingService, err := NewEmbeddingService(agent.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create embedding service: %w", err)
	}

	collectionName := k.getCollectionName(agentID)

	// Разбиваем документ на чанки
	metadata := map[string]interface{}{
		"agent_id":  agentID.String(),
		"file_name": fileName,
		"file_type": fileType,
		"file_size": fileSize,
		"added_at":  time.Now().UTC(),
	}

	chunks, err := embeddingService.SplitDocument(content, fileName, metadata)
	if err != nil {
		return fmt.Errorf("failed to split document: %w", err)
	}

	if len(chunks) == 0 {
		return fmt.Errorf("document produced no chunks")
	}

	// Генерируем эмбеддинги
	embeddings, err := embeddingService.GenerateEmbeddings(chunks)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Подготавливаем векторные точки
	points := make([]database.VectorPoint, len(chunks))
	for i, chunk := range chunks {
		points[i] = database.VectorPoint{
			ID:      chunk.ID,
			Vector:  embeddings[i],
			Payload: chunk.Metadata,
		}
	}

	// Сохраняем в векторную базу данных
	err = k.vectorDB.UpsertPoints(collectionName, points)
	if err != nil {
		return fmt.Errorf("failed to store embeddings: %w", err)
	}

	// Сохраняем метаданные файла в SQL базе данных
	knowledgeFile := &models.KnowledgeFile{
		AgentID:        agentID,
		FileName:       fileName,
		OriginalName:   fileName,
		FileSize:       fileSize,
		FileType:       fileType,
		CollectionName: collectionName,
		ChunkCount:     len(chunks),
		Status:         "completed",
		CreatedAt:      time.Now(),
		ProcessedAt:    &time.Time{},
	}
	*knowledgeFile.ProcessedAt = time.Now()

	err = k.db.CreateKnowledgeFile(knowledgeFile)
	if err != nil {
		return fmt.Errorf("failed to store file metadata: %w", err)
	}

	return nil
}

func (k *KnowledgeService) SearchKnowledge(agentID uuid.UUID, query string, limit int, scoreThreshold float32) ([]KnowledgeSearchResult, error) {
	// Получаем агента для использования его API ключа
	agent, err := k.db.GetAgentByID(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("agent not found: %w", err)
	}

	// Создаем сервис эмбеддингов с API ключом агента
	embeddingService, err := NewEmbeddingService(agent.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding service: %w", err)
	}

	collectionName := k.getCollectionName(agentID)

	// Генерируем эмбеддинг для запроса
	queryVector, err := embeddingService.GenerateQueryEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Ищем в векторной базе данных
	results, err := k.vectorDB.Search(collectionName, queryVector, limit, scoreThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to search knowledge base: %w", err)
	}

	// Конвертируем в результаты поиска знаний
	knowledgeResults := make([]KnowledgeSearchResult, len(results))
	for i, result := range results {
		content := ""
		source := ""

		if val, ok := result.Payload["content"].(string); ok {
			content = val
		}
		if val, ok := result.Payload["source"].(string); ok {
			source = val
		}

		knowledgeResults[i] = KnowledgeSearchResult{
			Content:  content,
			Source:   source,
			Score:    result.Score,
			Metadata: result.Payload,
		}
	}

	return knowledgeResults, nil
}

func (k *KnowledgeService) GetKnowledgeFiles(agentID uuid.UUID) ([]models.KnowledgeFile, error) {
	return k.db.GetKnowledgeFiles(agentID.String())
}

func (k *KnowledgeService) DeleteKnowledgeFile(agentID uuid.UUID, fileID uint) error {
	// Получаем информацию о файле
	file, err := k.db.GetKnowledgeFileByID(fileID)
	if err != nil {
		return fmt.Errorf("knowledge file not found: %w", err)
	}

	// Проверяем принадлежность файла агенту
	if file.AgentID != agentID {
		return fmt.Errorf("file does not belong to this agent")
	}

	// Получаем имя коллекции
	collectionName := k.getCollectionName(agentID)

	// Создаем фильтр для удаления векторов, связанных с этим файлом
	filter := map[string]interface{}{
		"file_name": file.FileName,
	}

	// Удаляем векторы из Qdrant по метаданным файла
	err = k.vectorDB.DeletePointsByFilter(collectionName, filter)
	if err != nil {
		return fmt.Errorf("failed to delete vectors from vector database: %w", err)
	}

	// Удаляем метаданные файла из SQL базы данных
	err = k.db.DeleteKnowledgeFile(fileID)
	if err != nil {
		return fmt.Errorf("failed to delete knowledge file: %w", err)
	}

	return nil
}

func (k *KnowledgeService) DeleteAllAgentKnowledge(agentID uuid.UUID) error {
	collectionName := k.getCollectionName(agentID)

	// Удаляем коллекцию из векторной базы данных
	err := k.vectorDB.DeleteCollection(collectionName)
	if err != nil {
		return fmt.Errorf("failed to delete vector collection: %w", err)
	}

	// Удаляем метаданные файлов из SQL базы данных
	err = k.db.DeleteKnowledgeFilesByAgent(agentID.String())
	if err != nil {
		return fmt.Errorf("failed to delete knowledge file metadata: %w", err)
	}

	return nil
}

func (k *KnowledgeService) GetKnowledgeStats(agentID uuid.UUID) (map[string]interface{}, error) {
	collectionName := k.getCollectionName(agentID)

	// Получаем количество векторов
	vectorCount, err := k.vectorDB.CountPoints(collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector count: %w", err)
	}

	// Получаем информацию о коллекции
	collectionInfo, err := k.vectorDB.GetCollectionInfo(collectionName)
	if err != nil {
		collectionInfo = map[string]interface{}{}
	}

	// Получаем количество файлов
	files, err := k.db.GetKnowledgeFiles(agentID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge files: %w", err)
	}

	totalSize := int64(0)
	totalChunks := 0
	for _, file := range files {
		totalSize += file.FileSize
		totalChunks += file.ChunkCount
	}

	stats := map[string]interface{}{
		"vector_count":    vectorCount,
		"file_count":      len(files),
		"total_size":      totalSize,
		"total_chunks":    totalChunks,
		"collection":      collectionName,
		"collection_info": collectionInfo,
	}

	return stats, nil
}

func (k *KnowledgeService) ValidateAgentKnowledge(agentID uuid.UUID) error {
	// Получаем агента
	agent, err := k.db.GetAgentByID(agentID.String())
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	// Проверяем API ключ агента
	embeddingService, err := NewEmbeddingService(agent.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create embedding service: %w", err)
	}

	err = embeddingService.ValidateAPIKey()
	if err != nil {
		return fmt.Errorf("invalid API key for agent: %w", err)
	}

	// Проверяем существование коллекции
	collectionName := k.getCollectionName(agentID)
	exists, err := k.vectorDB.CollectionExists(collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("knowledge collection does not exist")
	}

	return nil
}

func (k *KnowledgeService) getCollectionName(agentID uuid.UUID) string {
	return fmt.Sprintf("agent_%s", agentID.String())
}
