package api

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"io"
	"macdent-ai-chatbot/v2/configs"
	"macdent-ai-chatbot/v2/databases"
	"macdent-ai-chatbot/v2/services/agent"
	"macdent-ai-chatbot/v2/services/dialog"
	"macdent-ai-chatbot/v2/services/knowledge"
	"macdent-ai-chatbot/v2/utils"
	"mime/multipart"
)

type AgentHandler struct {
	config    *configs.ApiServerConfig
	postgres  *databases.PostgresDatabase
	qdrant    *databases.QdrantDatabase
	validator *validator.Validate
	loggger   *log.Logger
}

func NewAgentHandler(config *configs.ApiServerConfig) *AgentHandler {
	postgres := databases.NewPostgres(config.Postgres)
	qdrant := databases.NewQdrant(config.Qdrant)
	logger := utils.NewLogger("handler")

	return &AgentHandler{
		config:    config,
		postgres:  postgres,
		qdrant:    qdrant,
		validator: validator.New(),
		loggger:   logger,
	}
}

func (h *AgentHandler) CreateAgent(c fiber.Ctx) error {
	var request agent.CreateAgentRequest

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ошибка": "Неправильное тело запроса",
			"детали": err.Error(),
		})
	}

	err := h.validator.Struct(&request)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil && len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": newAgent,
	})
}

func (h *AgentHandler) GetAgents(c fiber.Ctx) error {
	var request agent.GetAgentsRequest
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ошибка": "Неправильное тело запроса",
			"детали": err.Error(),
		})
	}

	err := h.validator.Struct(&request)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil && len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": agents,
	})
}

func (h *AgentHandler) GetAgent(c fiber.Ctx) error {
	agentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": agents,
	})
}

func (h *AgentHandler) UpdateAgent(c fiber.Ctx) error {
	agentID := c.Params("id")

	var request agent.UpdateAgentRequest
	request.AgentID = agentID

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ошибка": "Неправильное тело запроса",
			"детали": err.Error(),
		})
	}

	err := h.validator.Struct(&request)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil && len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
	agentID := c.Params("id")

	var request knowledge.UploadKnowledgeRequest
	request.AgentID = agentID

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ошибка": "Ошибка парсинга multipart формы",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ошибка": "Не переданы файлы для загрузки",
		})
	}

	for _, file := range files {
		if file.Size == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ошибка": "Файл не должен быть пустым",
			})
		}

		err = func() error {
			fileStream, err := file.Open()
			if err != nil {
				return err
			}
			defer func(fileStream multipart.File) {
				if closeErr := fileStream.Close(); closeErr != nil {
					h.loggger.Errorf("закрытие файла: %v", closeErr)
				}
			}(fileStream)

			content, err := io.ReadAll(fileStream)
			if err != nil {
				return err
			}

			knowledgeFile := knowledge.Knowledge{
				Name:    file.Filename,
				Size:    file.Size,
				Type:    file.Header.Get("Content-Type"),
				Content: content,
			}

			request.Files = append(request.Files, knowledgeFile)
			return nil
		}()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ошибка": "Ошибка обработки файла",
				"детали": err.Error(),
			})
		}
	}

	if err := h.validator.Struct(&request); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ошибка": "Неправильные данные запроса",
				"детали": validationErrors.Error(),
			})
		}
	}

	errorResponse := knowledge.NewService(h.postgres, h.qdrant).
		UploadKnowledge(&request)

	if errorResponse != nil {
		return c.Status(errorResponse.StatusCode).JSON(fiber.Map{
			"ошибка": errorResponse.Message,
			"детали": errorResponse.Details,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"сообщение": "Файлы успешно загружены",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ошибка": "Неправильное тело запроса",
			"детали": err.Error(),
		})
	}

	err := h.validator.Struct(&request)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	if err != nil && len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
