package tool

import "github.com/openai/openai-go"

func (s *Service) GetToolsFunctions() []openai.ChatCompletionToolParam {
	s.logger.Info("получение списка функций инструментов")

	return []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F("get_doctors"),
				Description: openai.String("Получает список врачей с данными"),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]string{
							"type": "string",
						},
					},
				}),
			}),
		},
	}
}
