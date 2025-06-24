package api

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"macdent-ai-chatbot/v2/configs"
	"macdent-ai-chatbot/v2/utils"
	"strconv"
)

type Server struct {
	app    *fiber.App
	config *configs.ApiServerConfig
	logger *log.Logger
}

func NewServer(cfg *configs.ApiServerConfig) *Server {
	customLogger := utils.NewLogger("api")

	app := fiber.New(fiber.Config{
		AppName:      "MacDent AI",
		ErrorHandler: errorHandler,
	})

	app.Use(logger.New(logger.Config{
		Done: func(c fiber.Ctx, logString []byte) {
			customLogger.Info(string(logString))
		},
	}))

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	return &Server{
		app:    app,
		config: cfg,
		logger: customLogger,
	}
}

func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Неизвестная ошибка"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"ошибка": message,
		"путь":   c.Path(),
		"метод":  c.Method(),
	})
}

func (s *Server) Setup() {
	api := s.app.Group("/api/v1")

	agentHandler := NewAgentHandler(s.config)
	agents := api.Group("/agents")

	// Создание агента
	agents.Post("/", agentHandler.CreateAgent)
	// Получение списка агентов
	agents.Get("/", agentHandler.GetAgents)

	// Получение конкретного агента
	agents.Get("/:id", agentHandler.GetAgent)
	// Обновление агента
	agents.Patch("/:id", agentHandler.UpdateAgent)
	// Удаление агента
	agents.Delete("/:id", agentHandler.DeleteAgent)

	// Загрузка базы знаний
	agents.Post("/:id/knowledge", agentHandler.UploadKnowledge)
	// Получение базы знаний
	agents.Get("/:id/knowledge", agentHandler.GetKnowledge)
	// Удаление базы знаний
	agents.Delete("/:id/knowledge", agentHandler.DeleteKnowledge)

	// Получение диалогов агента
	agents.Get("/:id/dialogs", agentHandler.GetDialogs)
	// Запрос на ответ диалогу
	agents.Post("/:id/dialogs", agentHandler.ResponseDialog)

	openaiHandler := NewOpenAIHandler()
	openai := api.Group("/openai")

	// Получение доступных моделей для агента
	openai.Get("/models", openaiHandler.GetModels)

	mockHandler := NewMockHandler()
	mock := api.Group("/mocks")

	// Получение списка докторов
	mock.Get("/doctors", mockHandler.GetDoctors)
}

func (s *Server) Run() {

	port := strconv.Itoa(s.config.Port)

	s.logger.Printf("🚀 Сервер запущен на порту: %s", port)

	err := s.app.Listen(":" + port)
	if err != nil {
		s.logger.Fatalf("API сервер: %v", err)
	}
}
