package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func ValidateRequest(c fiber.Ctx, request interface{}) error {
	if err := c.Bind().JSON(request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Неправильное тело запроса",
			"детали": err.Error(),
		})
	}

	err := validate.Struct(request)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil && len(validationErrors) > 0 {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Неправильные данные запроса",
			"детали": validationErrors.Error(),
		})
	}

	return nil
}
