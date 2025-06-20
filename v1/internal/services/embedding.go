package services

import (
	"context"
	"fmt"
	"hash/fnv"
	"macdent-ai-chatbot/v1/internal/utlis"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type EmbeddingService struct {
	client     *openai.Client
	splitter   *utils.utils
	vectorSize int
	model      string
}

type DocumentChunk struct {
	ID       uint64
	Content  string
	Metadata map[string]interface{}
}

func NewEmbeddingService(apiKey string) (*EmbeddingService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	splitter := utils.NewRecursiveTextSplitter()

	return &EmbeddingService{
		client:     client,
		splitter:   splitter,
		vectorSize: 1536,
		model:      "text-embedding-3-small",
	}, nil
}

func (e *EmbeddingService) GetVectorSize() int {
	return e.vectorSize
}

func (e *EmbeddingService) SplitDocument(content string, source string, metadata map[string]interface{}) ([]DocumentChunk, error) {
	// Очищаем текст
	content = utils.CleanText(content)

	if !utils.IsValidUTF8(content) {
		return nil, fmt.Errorf("invalid UTF-8 content in document")
	}

	chunks := e.splitter.SplitText(content)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("document produced no chunks")
	}

	documentChunks := make([]DocumentChunk, len(chunks))
	for i, chunk := range chunks {
		chunkMetadata := make(map[string]interface{})

		// Копируем исходные метаданные
		for k, v := range metadata {
			chunkMetadata[k] = v
		}

		// Добавляем метаданные чанка
		chunkMetadata["source"] = source
		chunkMetadata["chunk_index"] = i
		chunkMetadata["content"] = chunk
		chunkMetadata["chunk_size"] = len(chunk)
		chunkMetadata["token_count"] = utils.CountTokensApprox(chunk)

		documentChunks[i] = DocumentChunk{
			ID:       e.generateChunkID(chunk, source, i),
			Content:  chunk,
			Metadata: chunkMetadata,
		}
	}

	return documentChunks, nil
}

func (e *EmbeddingService) GenerateEmbeddings(chunks []DocumentChunk) ([][]float64, error) {
	if len(chunks) == 0 {
		return [][]float64{}, nil
	}

	// Подготавливаем тексты для эмбеддинга
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}

	// Генерируем эмбеддинги
	resp, err := e.client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Model: openai.F(e.model),
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](
			openai.EmbeddingNewParamsInputArrayOfStrings(texts),
		),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	if len(resp.Data) != len(chunks) {
		return nil, fmt.Errorf("embedding count mismatch: expected %d, got %d", len(chunks), len(resp.Data))
	}

	// Конвертируем результат
	results := make([][]float64, len(resp.Data))
	for i, embedding := range resp.Data {
		results[i] = embedding.Embedding
	}

	return results, nil
}

func (e *EmbeddingService) GenerateQueryEmbedding(query string) ([]float64, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	resp, err := e.client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Model: openai.F(e.model),
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](
			openai.EmbeddingNewParamsInputArrayOfStrings([]string{query}),
		),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned for query")
	}

	return resp.Data[0].Embedding, nil
}

func (e *EmbeddingService) ValidateAPIKey() error {
	// Проверяем API ключ простым запросом
	_, err := e.client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Model: openai.F(e.model),
		Input: openai.F[openai.EmbeddingNewParamsInputUnion](
			openai.EmbeddingNewParamsInputArrayOfStrings([]string{"test"}),
		),
	})

	return err
}

func (e *EmbeddingService) generateChunkID(content, source string, index int) uint64 {
	h := fnv.New64a()
	_, err := h.Write([]byte(fmt.Sprintf("%s:%s:%d", source, content, index)))
	if err != nil {
		fmt.Printf("error generating chunk ID: %v\n", err)
		os.Exit(1)
	}
	return h.Sum64()
}
