package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"macdent-ai-chatbot/v2/services/openai"
)

type OpenAIHandler struct {
	validator *validator.Validate
}

func NewOpenAIHandler() *OpenAIHandler {
	return &OpenAIHandler{
		validator: validator.New(),
	}
}

func (h *OpenAIHandler) GetModels(c fiber.Ctx) error {
	var request openai.GetModelsRequest

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Неправильное тело запроса",
			"детали": err.Error(),
		})
	}

	err := h.validator.Struct(&request)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil && len(validationErrors) > 0 {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Неправильные данные запроса",
			"детали": validationErrors.Error(),
		})
	}

	models, errorResponse := openai.NewService(request.ApiKey).
		GetModels()

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": models,
	})
}
