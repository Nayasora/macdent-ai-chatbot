package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"macdent-ai-chatbot/v2/configs"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/services/agent"
)

type AgentHandler struct {
	config   *configs.ApiServerConfig
	postgres *databases.PostgresDatabase
	qdrant   *databases.QdrantDatabase
}

func NewAgentHandler(config *configs.ApiServerConfig) *AgentHandler {
	postgres := databases.NewPostgres(config.Postgres)
	qdrant := databases.NewQdrant(config.Qdrant)

	return &AgentHandler{
		config:   config,
		postgres: postgres,
		qdrant:   qdrant,
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

	newAgent := agent.NewService().
		CreateAgent(&request, h.postgres)

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

	agents := agent.NewService().
		GetAgents(&request, h.postgres)

	return c.Status(200).JSON(fiber.Map{
		"data": agents,
	})
}

func (h *AgentHandler) GetAgent(c fiber.Ctx) error {
	// TODO: Реализовать получение конкретного агента
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
	})
}

func (h *AgentHandler) UpdateAgent(c fiber.Ctx) error {
	// TODO: Реализовать обновление агента
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
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

	agent.NewService().
		GetAgent(
			&agent.GetAgentRequest{
				AgentID: agentID.String(),
			},
			h.postgres,
		)

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
	// TODO: Реализовать запрос на ответ в диалоге
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"сообщение": "Метод не реализован",
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
