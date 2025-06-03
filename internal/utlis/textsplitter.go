package utils

import (
	"strings"
	"unicode"
)

type TextSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string
}

type DocumentChunk struct {
	Text     string
	Metadata map[string]interface{}
}

func NewRecursiveTextSplitter() *TextSplitter {
	return &TextSplitter{
		ChunkSize:    512,
		ChunkOverlap: 100,
		Separators:   []string{"\n\n", "\n", " ", ""},
	}
}

func (ts *TextSplitter) SplitText(text string) []string {
	if len(text) <= ts.ChunkSize {
		return []string{text}
	}

	return ts.splitTextRecursive(text, ts.Separators)
}

func (ts *TextSplitter) splitTextRecursive(text string, separators []string) []string {
	var finalChunks []string
	separator := separators[len(separators)-1]

	var newSeparators []string
	for _, sep := range separators {
		if sep == "" {
			separator = sep
			break
		}
		if strings.Contains(text, sep) {
			separator = sep
			newSeparators = separators[indexOf(separators, sep)+1:]
			break
		}
	}

	splits := ts.splitTextWithSeparator(text, separator)

	var goodSplits []string
	mergeSeparator := ""
	if separator != "" {
		mergeSeparator = separator
	}

	for _, split := range splits {
		if len(split) < ts.ChunkSize {
			goodSplits = append(goodSplits, split)
		} else {
			if len(goodSplits) > 0 {
				mergedText := ts.mergeSplits(goodSplits, mergeSeparator)
				finalChunks = append(finalChunks, mergedText...)
				goodSplits = []string{}
			}
			if len(newSeparators) == 0 {
				finalChunks = append(finalChunks, split)
			} else {
				otherInfo := ts.splitTextRecursive(split, newSeparators)
				finalChunks = append(finalChunks, otherInfo...)
			}
		}
	}

	if len(goodSplits) > 0 {
		mergedText := ts.mergeSplits(goodSplits, mergeSeparator)
		finalChunks = append(finalChunks, mergedText...)
	}

	return finalChunks
}

func (ts *TextSplitter) splitTextWithSeparator(text, separator string) []string {
	if separator == "" {
		return ts.splitByCharacter(text)
	}
	return strings.Split(text, separator)
}

func (ts *TextSplitter) splitByCharacter(text string) []string {
	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += ts.ChunkSize {
		end := i + ts.ChunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}

	return chunks
}

func (ts *TextSplitter) mergeSplits(splits []string, separator string) []string {
	var docs []string
	var currentDoc []string
	total := 0

	for _, split := range splits {
		length := len(split)
		if total+length+(len(currentDoc)*len(separator)) > ts.ChunkSize && len(currentDoc) > 0 {
			if len(currentDoc) > 0 {
				doc := strings.Join(currentDoc, separator)
				if strings.TrimSpace(doc) != "" {
					docs = append(docs, doc)
				}

				// Handle overlap
				for total > ts.ChunkOverlap || (total+length+(len(currentDoc)*len(separator)) > ts.ChunkSize && total > 0) {
					if len(currentDoc) == 0 {
						break
					}
					removed := currentDoc[0]
					currentDoc = currentDoc[1:]
					total -= len(removed) + len(separator)
				}
			}
		}
		currentDoc = append(currentDoc, split)
		total += length + len(separator)
	}

	if len(currentDoc) > 0 {
		doc := strings.Join(currentDoc, separator)
		if strings.TrimSpace(doc) != "" {
			docs = append(docs, doc)
		}
	}

	return docs
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

// Дополнительные утилиты для обработки текста
func CleanText(text string) string {
	// Удаляем лишние пробелы и переносы строк
	text = strings.TrimSpace(text)

	// Заменяем множественные пробелы на одинарные
	words := strings.Fields(text)
	return strings.Join(words, " ")
}

func CountTokensApprox(text string) int {
	// Приблизительный подсчет токенов (4 символа ≈ 1 токен для английского)
	// Для более точного подсчета можно использовать tiktoken
	runes := []rune(text)
	return len(runes) / 4
}

func IsValidUTF8(text string) bool {
	for _, r := range text {
		if r == unicode.ReplacementChar {
			return false
		}
	}
	return true
}

func TruncateText(text string, maxLength int) string {
	runes := []rune(text)
	if len(runes) <= maxLength {
		return text
	}
	return string(runes[:maxLength]) + "..."
}
