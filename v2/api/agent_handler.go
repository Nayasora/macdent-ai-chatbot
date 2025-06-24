package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"macdent-ai-chatbot/v2/configs"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/services/agent"
	"macdent-ai-chatbot/v2/services/dialog"
)

type AgentHandler struct {
	config    *configs.ApiServerConfig
	postgres  *databases.PostgresDatabase
	qdrant    *databases.QdrantDatabase
	validator *validator.Validate
}

func NewAgentHandler(config *configs.ApiServerConfig) *AgentHandler {
	postgres := databases.NewPostgres(config.Postgres)
	qdrant := databases.NewQdrant(config.Qdrant)

	return &AgentHandler{
		config:    config,
		postgres:  postgres,
		qdrant:    qdrant,
		validator: validator.New(),
	}
}

func (h *AgentHandler) CreateAgent(c fiber.Ctx) error {
	var request agent.CreateAgentRequest

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

	newAgent, errorResponse := agent.NewService().
		CreateAgent(&request, h.postgres)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": newAgent,
	})
}

func (h *AgentHandler) GetAgents(c fiber.Ctx) error {
	var request agent.GetAgentsRequest
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

	agents, errorResponse := agent.NewService().
		GetAgents(&request, h.postgres)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": agents,
	})
}

func (h *AgentHandler) GetAgent(c fiber.Ctx) error {
	agentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Неверный формат ID агента",
		})
	}

	agents, errorResponse := agent.NewService().
		GetAgent(agentID, h.postgres)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": agents,
	})
}

func (h *AgentHandler) UpdateAgent(c fiber.Ctx) error {
	agentID := c.Params("id")

	var request agent.UpdateAgentRequest
	request.AgentID = agentID

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

	updatedAgent, errorResponse := agent.NewService().
		UpdateAgent(&request, h.postgres)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedAgent,
	})
}

func (h *AgentHandler) DeleteAgent(c fiber.Ctx) error {
	// TODO: Реализовать удаление агента
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) UploadKnowledge(c fiber.Ctx) error {
	agentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Неверный формат ID агента",
		})
	}

	_, errorResponse := agent.NewService().
		GetAgent(
			agentID,
			h.postgres,
		)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Ошибка парсинга multipart формы",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"ошибка": "Не переданы файлы для загрузки",
		})
	}

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) GetKnowledge(c fiber.Ctx) error {
	// TODO: Реализовать получение базы знаний
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) DeleteKnowledge(c fiber.Ctx) error {
	// TODO: Реализовать удаление базы знаний
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) GetDialogs(c fiber.Ctx) error {
	// TODO: Реализовать получение диалогов агента
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) DeleteDialogs(c fiber.Ctx) error {
	// TODO: Реализовать удаление всех диалогов агента
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) CreateDialog(c fiber.Ctx) error {
	// TODO: Реализовать создание диалога
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) ResponseDialog(c fiber.Ctx) error {
	agentID := c.Params("id")

	var request dialog.UserDialogNewMessageRequest
	request.AgentID = agentID

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

	agentMessage, errorResponse := dialog.NewService().
		ResponseDialogNewMessageRequest(&request, h.postgres)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": agentMessage,
	})
}

func (h *AgentHandler) GetDialog(c fiber.Ctx) error {
	// TODO: Реализовать получение конкретного диалога
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) DeleteDialog(c fiber.Ctx) error {
	// TODO: Реализовать удаление конкретного диалога
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}
