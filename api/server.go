package api

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"
	"macdent-ai-chatbot/config"
	"macdent-ai-chatbot/logger"
	"strconv"
)

type Server struct {
	App    *fiber.App
	Logger *log.Logger
	Config *config.ApiServerConfig
}

func NewServer(cfg *config.ApiServerConfig) *Server {
	customLogger := logger.New("api")
	fiberServer := fiber.New()

	return &Server{
		App:    fiberServer,
		Logger: customLogger,
		Config: cfg,
	}
}

func (s *Server) Run() {
	s.App.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	err := s.App.Listen(":" + strconv.Itoa(s.Config.Port))
	if err != nil {
		s.Logger.Fatalf("запуск сервера: %s", err)
	}
}
