package openai

import (
	"strings"
)

func (s *Service) SplitSentences(text string) []string {
	text = strings.ReplaceAll(text, "\n\n", ". ")
	text = strings.ReplaceAll(text, "\n", ". ")

	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?' || r == ';'
	})

	s.logger.Infof("количество предложений: %d", len(sentences))

	var cleanSentences []string
	for _, sentence := range sentences {
		cleaned := strings.TrimSpace(sentence)
		if len(cleaned) > 0 {
			cleanSentences = append(cleanSentences, cleaned)
			s.logger.Infof("новое предложение: %s", cleaned)
		}
	}

	return cleanSentences
}
