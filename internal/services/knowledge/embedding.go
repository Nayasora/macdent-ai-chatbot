package knowledge

import (
	"context"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/qdrant/go-client/qdrant"
	openaiService "macdent-ai-chatbot/internal/services/openai"
	"macdent-ai-chatbot/internal/utils"
	"strings"
	"time"
	"unicode/utf8"
)

type EmbeddingResult struct {
	Chunk      Chunk
	Vector     []float32
	TokenUsage int
}

const embeddingsBatchSize = 50

func (s *Service) PrepareContent(openaiService *openaiService.Service, content string, agentID uuid.UUID) ([]EmbeddingResult, *utils.UserErrorResponse) {
	if strings.TrimSpace(content) == "" {
		return nil, utils.NewUserErrorResponse(400, "Пустой контент", "Контент не может быть пустым")
	}

	if !utf8.ValidString(content) {
		s.logger.Errorf("входной контент содержит невалидные UTF-8 символы")
		return nil, utils.NewUserErrorResponse(400, "Невалидный контент", "Контент содержит недопустимые символы")
	}

	config := NewDefaultChunkingConfig()
	chunks := s.CreateChunks(content, config)

	if len(chunks) == 0 {
		return nil, utils.NewUserErrorResponse(400, "Не удалось создать чанки", "Проверьте корректность контента")
	}

	s.logger.Infof("создано %d чанков для агента %s", len(chunks), agentID.String())

	validChunks := make([]Chunk, 0, len(chunks))
	for i, chunk := range chunks {
		if !utf8.ValidString(chunk.Text) {
			s.logger.Warnf("пропускаем невалидный чанк %d", i)
			continue
		}
		if strings.TrimSpace(chunk.Text) == "" {
			s.logger.Warnf("пропускаем пустой чанк %d", i)
			continue
		}
		validChunks = append(validChunks, chunk)
	}

	if len(validChunks) == 0 {
		return nil, utils.NewUserErrorResponse(400, "Нет валидных чанков", "Все чанки содержат ошибки")
	}

	var results []EmbeddingResult

	for i := 0; i < len(validChunks); i += embeddingsBatchSize {
		end := i + embeddingsBatchSize
		if end > len(validChunks) {
			end = len(validChunks)
		}

		batch := validChunks[i:end]
		batchResults, err := s.processChunkBatch(openaiService, batch, agentID)
		if err != nil {
			return nil, err
		}

		results = append(results, batchResults...)

		if i+embeddingsBatchSize < len(validChunks) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return results, nil
}

func (s *Service) UpsertChunks(ctx context.Context, agentID uuid.UUID, results []EmbeddingResult) *utils.UserErrorResponse {
	if len(results) == 0 {
		return nil
	}

	points := make([]*qdrant.PointStruct, len(results))
	for i, result := range results {
		points[i] = &qdrant.PointStruct{
			Id:      qdrant.NewID(uuid.NewString()),
			Vectors: qdrant.NewVectors(result.Vector...),
			Payload: qdrant.NewValueMap(map[string]any{
				"text":        result.Chunk.Text,
				"chunk_index": int64(result.Chunk.Metadata["chunk_index"].(int)),
				"char_count":  int64(result.Chunk.Metadata["char_count"].(int)),
				"word_count":  int64(result.Chunk.Metadata["word_count"].(int)),
				"start_idx":   int64(result.Chunk.StartIdx),
				"end_idx":     int64(result.Chunk.EndIdx),
				"agent_id":    agentID.String(),
				"created_at":  time.Now().Format(time.RFC3339),
			}),
		}
	}

	return s.UpsertPoints(ctx, agentID, points)
}

func (s *Service) processChunkBatch(openaiService *openaiService.Service, chunks []Chunk, agentID uuid.UUID) ([]EmbeddingResult, *utils.UserErrorResponse) {
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Text
		s.logger.Infof("обработка чанка %d (%d символов): %s",
			chunk.Metadata["chunk_index"],
			len([]rune(chunk.Text)),
			s.truncateForLog(chunk.Text, 100))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := openaiService.Client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.String(openai.EmbeddingModelTextEmbedding3Large),
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](
			openai.EmbeddingNewParamsInputArrayOfStrings(texts),
		),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})

	if err != nil {
		s.logger.Errorf("создания эмбеддингов: %v", err)
		return nil, utils.NewUserErrorResponse(500, "Ошибка создания эмбеддингов", "Попробуйте позже")
	}

	if len(response.Data) != len(chunks) {
		s.logger.Errorf("несоответствие количества эмбеддингов: ожидалось %d, получено %d", len(chunks), len(response.Data))
		return nil, utils.NewUserErrorResponse(500, "Ошибка обработки эмбеддингов", "Внутренняя ошибка")
	}

	results := make([]EmbeddingResult, len(chunks))
	for i, embedding := range response.Data {
		vector := make([]float32, len(embedding.Embedding))
		for j, v := range embedding.Embedding {
			vector[j] = float32(v)
		}

		results[i] = EmbeddingResult{
			Chunk:      chunks[i],
			Vector:     vector,
			TokenUsage: int(response.Usage.TotalTokens) / len(response.Data),
		}
	}

	s.logger.Infof("успешно создано %d эмбеддингов для агента %s", len(results), agentID.String())
	return results, nil
}

func (s *Service) truncateForLog(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	return string(runes[:maxLen]) + "..."
}
