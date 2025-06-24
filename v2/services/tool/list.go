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
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F("get_schedule"),
				Description: openai.String("Получает расписание врачей"),
				Parameters: openai.F(openai.FunctionParameters{
					"type":       "object",
					"properties": map[string]interface{}{},
				}),
			}),
		},
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F("create_patient"),
				Description: openai.String("Создает пациента"),
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
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F("create_appointment"),
				Description: openai.String("Создает запись к врачу"),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"patient": map[string]string{
							"type":        "integer",
							"description": "ID пациента, полученный после вызова create_patient",
						},
						"doctor": map[string]string{
							"type":        "integer",
							"description": "ID врача, полученный из списка врачей (get_doctors)",
						},
						"date": map[string]string{
							"type": "string",
						},
						"start": map[string]string{
							"type": "string",
						},
						"end": map[string]string{
							"type": "string",
						},
					},
				}),
			}),
		},
	}
}
