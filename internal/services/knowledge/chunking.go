package knowledge

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type ChunkingConfig struct {
	ChunkSize    int
	OverlapSize  int
	MinChunkSize int
}

type Chunk struct {
	Text     string
	StartIdx int
	EndIdx   int
	Metadata map[string]interface{}
}

func NewDefaultChunkingConfig() *ChunkingConfig {
	return &ChunkingConfig{
		ChunkSize:    800, // Уменьшил для более стабильной работы
		OverlapSize:  150,
		MinChunkSize: 50,
	}
}

func (s *Service) CreateChunks(text string, config *ChunkingConfig) []Chunk {
	if config == nil {
		config = NewDefaultChunkingConfig()
	}

	text = s.preprocessText(text)
	runes := []rune(text) // Работаем с рунами, а не байтами!

	if len(runes) <= config.ChunkSize {
		return []Chunk{{
			Text:     text,
			StartIdx: 0,
			EndIdx:   len(runes),
			Metadata: map[string]interface{}{"chunk_index": 0},
		}}
	}

	var chunks []Chunk
	start := 0
	chunkIndex := 0

	for start < len(runes) {
		end := start + config.ChunkSize

		if end >= len(runes) {
			end = len(runes)
		} else {
			// Ищем оптимальную точку разбиения
			end = s.findOptimalBreakpointRunes(runes, start, end)
		}

		chunkText := strings.TrimSpace(string(runes[start:end]))

		// Проверяем валидность UTF-8
		if !utf8.ValidString(chunkText) {
			s.logger.Errorf("невалидная UTF-8 строка в чанке %d", chunkIndex)
			continue
		}

		if len([]rune(chunkText)) >= config.MinChunkSize {
			chunks = append(chunks, Chunk{
				Text:     chunkText,
				StartIdx: start,
				EndIdx:   end,
				Metadata: map[string]interface{}{
					"chunk_index": chunkIndex,
					"char_count":  len([]rune(chunkText)),
					"word_count":  len(strings.Fields(chunkText)),
				},
			})
			chunkIndex++
		}

		// Вычисляем следующую позицию с учетом оверлапа
		nextStart := end - config.OverlapSize
		if nextStart <= start {
			nextStart = start + (config.ChunkSize - config.OverlapSize)
		}

		if nextStart >= len(runes) {
			break
		}

		start = nextStart
	}

	return chunks
}

func (s *Service) preprocessText(text string) string {
	// Нормализуем переносы строк
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Убираем лишние пробелы и пустые строки
	lines := strings.Split(text, "\n")
	var processedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Заменяем множественные пробелы на одиночные
			line = strings.Join(strings.Fields(line), " ")
			processedLines = append(processedLines, line)
		}
	}

	return strings.Join(processedLines, " ")
}

func (s *Service) findOptimalBreakpointRunes(runes []rune, start, maxEnd int) int {
	if maxEnd >= len(runes) {
		return len(runes)
	}

	// Ищем оптимальную точку разбиения в пределах последних 100 символов
	searchStart := maxEnd - 100
	if searchStart < start {
		searchStart = start
	}

	bestBreak := maxEnd

	// Ищем конец предложения
	for i := maxEnd - 1; i >= searchStart; i-- {
		char := runes[i]

		if char == '.' || char == '!' || char == '?' {
			// Проверяем, что после знака препинания идет пробел или конец текста
			if i+1 < len(runes) && unicode.IsSpace(runes[i+1]) {
				return i + 1
			}
			if i+1 == len(runes) {
				return i + 1
			}
		}
	}

	// Если не нашли конец предложения, ищем пробел
	for i := maxEnd - 1; i >= searchStart; i-- {
		if unicode.IsSpace(runes[i]) {
			bestBreak = i
			break
		}
	}

	return bestBreak
}
