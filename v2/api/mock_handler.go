package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"macdent-ai-chatbot/v2/services/mock"
)

type MockHandler struct {
	validator *validator.Validate
}

func NewMockHandler() *MockHandler {
	return &MockHandler{
		validator: validator.New(),
	}
}

func (h *MockHandler) GetDoctors(c fiber.Ctx) error {
	var request mock.GetDoctorsRequest

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

	models := mock.NewService().
		GetDoctors(&request)

	return c.Status(200).JSON(fiber.Map{
		"data": models,
	})
}
