package api

import (
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

	if err := ValidateRequest(c, &request); err != nil {
		return err
	}

	models := openai.NewService(request.ApiKey).
		GetModels()

	return c.Status(200).JSON(fiber.Map{
		"data": models,
	})
}
