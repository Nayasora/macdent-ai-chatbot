package knowledge

import (
	"context"
	"github.com/openai/openai-go"
	openaiService "macdent-ai-chatbot/v2/services/openai"
	"macdent-ai-chatbot/v2/utils"
	"strings"
	"time"
)

type EmbeddingResult struct {
	Text       string
	Vector     []float64
	TokenUsage int
}

const batchSize = 100

func (s *Service) PrepareContent(openaiService *openaiService.Service, content string) ([]EmbeddingResult, *utils.UserErrorResponse) {
	if strings.TrimSpace(content) == "" {
		return nil, utils.NewUserErrorResponse(
			400,
			"Пустой контент",
			"Контент не может быть пустым",
		)
	}

	sentences := openaiService.SplitSentences(content)
	if len(sentences) == 0 {
		return nil, utils.NewUserErrorResponse(
			400,
			"Контент не содержит предложений",
			"Убедитесь, что контент содержит осмысленный текст",
		)
	}

	var results []EmbeddingResult

	for i := 0; i < len(sentences); i += batchSize {
		end := i + batchSize
		if end > len(sentences) {
			end = len(sentences)
		}

		batch := sentences[i:end]
		batchResults, err := s.processEmbeddingBatch(openaiService, batch)
		if err != nil {
			return nil, err
		}

		results = append(results, batchResults...)
	}

	return results, nil
}

func (s *Service) processEmbeddingBatch(openaiService *openaiService.Service, sentences []string) ([]EmbeddingResult, *utils.UserErrorResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := openaiService.Client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.String(openai.EmbeddingModelTextEmbedding3Large),
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](
			openai.EmbeddingNewParamsInputArrayOfStrings(sentences),
		),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})

	if err != nil {
		s.logger.Errorf("ошибка создания эмбеддингов: %v", err)
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка создания эмбеддингов",
			"Пожалуйста, попробуйте позже",
		)
	}

	if len(response.Data) != len(sentences) {
		s.logger.Errorf("несоответствие количества эмбеддингов: ожидалось %d, получено %d",
			len(sentences), len(response.Data))
		return nil, utils.NewUserErrorResponse(
			500,
			"Ошибка обработки эмбеддингов",
			"Внутренняя ошибка сервиса",
		)
	}

	results := make([]EmbeddingResult, len(sentences))
	for i, embedding := range response.Data {
		results[i] = EmbeddingResult{
			Text:       sentences[i],
			Vector:     embedding.Embedding,
			TokenUsage: int(response.Usage.TotalTokens) / len(response.Data),
		}
	}

	s.logger.Infof("успешно создано %d эмбеддингов", len(results))
	return results, nil
}
