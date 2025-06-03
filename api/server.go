package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"macdent-ai-chatbot/config"
	"macdent-ai-chatbot/internal/database"
	"macdent-ai-chatbot/internal/handlers"
	"macdent-ai-chatbot/internal/services"
)

type Server struct {
	app    *fiber.App
	config config.ApiServerConfig
}

func NewServer(cfg config.ApiServerConfig) *Server {
	app := fiber.New(fiber.Config{
		AppName:      "MacDent AI Chatbot",
		ErrorHandler: errorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	return &Server{
		app:    app,
		config: cfg,
	}
}

func (s *Server) SetupRoutes(cfg *config.Config) error {
	db, err := database.NewDatabase(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	vectorDB, err := database.NewVectorDatabase(cfg.Qdrant.Host, cfg.Qdrant.Port, cfg.Qdrant.ApiKey)
	if err != nil {
		return fmt.Errorf("failed to connect to Qdrant: %w", err)
	}

	knowledgeService := services.NewKnowledgeService(db, vectorDB)
	agentService := services.NewAgentService(db, knowledgeService)
	dialogService := services.NewDialogService(db, knowledgeService)

	agentHandler := handlers.NewAgentHandler(agentService, dialogService, knowledgeService)

	api := s.app.Group("/api/v1")

	agents := api.Group("/agents")
	agents.Post("/", agentHandler.CreateAgent)
	agents.Get("/", agentHandler.ListAgents)
	agents.Get("/:id", agentHandler.GetAgent)
	agents.Put("/:id", agentHandler.UpdateAgent)
	agents.Delete("/:id", agentHandler.DeleteAgent)
	agents.Get("/:id/stats", agentHandler.GetAgentStats)

	agents.Post("/:id/knowledge", agentHandler.UploadKnowledge)
	agents.Get("/:id/knowledge", agentHandler.GetKnowledge)
	agents.Delete("/:id/knowledge", agentHandler.DeleteKnowledge)

	agents.Post("/:id/dialog", agentHandler.ProcessDialog)
	agents.Get("/:id/dialogs", agentHandler.GetDialogHistory)

	s.app.Get("/health", func(c fiber.Ctx) error {
		if err := db.HealthCheck(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  "database connection failed",
			})
		}

		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": "2025-06-02 15:00:58",
			"version":   "1.0.0",
		})
	})

	s.app.Get("/api/docs", func(c fiber.Ctx) error {
		docs := map[string]interface{}{
			"title":       "MacDent AI Chatbot API",
			"version":     "1.0.0",
			"description": "API for managing AI agents with knowledge bases and dialog processing",
			"timestamp":   "2025-06-02 15:00:58",
			"endpoints": map[string]interface{}{
				"agents": map[string]string{
					"POST /api/v1/agents":          "Create new agent with individual OpenAI API key",
					"GET /api/v1/agents":           "List all active agents",
					"GET /api/v1/agents/:id":       "Get agent details by ID",
					"PUT /api/v1/agents/:id":       "Update agent configuration",
					"DELETE /api/v1/agents/:id":    "Delete agent and all associated data",
					"GET /api/v1/agents/:id/stats": "Get agent usage statistics",
				},
				"knowledge": map[string]string{
					"POST /api/v1/agents/:id/knowledge": "Upload knowledge files (supports multiple files)",
					"GET /api/v1/agents/:id/knowledge":  "Get knowledge base files and statistics",
				},
				"dialogs": map[string]string{
					"POST /api/v1/agents/:id/dialog": "Send message to agent (with RAG support)",
					"GET /api/v1/agents/:id/dialogs": "Get dialog history for specific user",
				},
				"system": map[string]string{
					"GET /health":   "System health check",
					"GET /api/docs": "This API documentation",
				},
			},
			"examples": map[string]interface{}{
				"create_agent": map[string]interface{}{
					"url": "POST /api/v1/agents",
					"body": map[string]interface{}{
						"name":             "Customer Support Bot",
						"description":      "AI assistant for customer support",
						"api_key":          "sk-...",
						"model":            "gpt-4o-mini",
						"context_size":     4096,
						"system_prompt":    "You are a helpful customer support assistant.",
						"assistant_prompt": "Always be polite and professional.",
						"temperature":      0.7,
						"max_tokens":       1000,
					},
				},
				"send_message": map[string]interface{}{
					"url": "POST /api/v1/agents/{agent_id}/dialog",
					"body": map[string]interface{}{
						"message":         "Hello, I need help with my order",
						"user_id":         "user123",
						"use_knowledge":   true,
						"knowledge_limit": 3,
						"score_threshold": 0.5,
						"history_limit":   10,
					},
				},
			},
		}
		return c.JSON(docs)
	})

	return nil
}

func (s *Server) Run() {
	cfg := config.NewConfig()

	if err := s.SetupRoutes(cfg); err != nil {
		log.Fatal("Failed to setup routes:", err)
	}

	log.Printf("🚀 Server starting on port %s", s.config.Port)
	log.Printf("📚 API Documentation: http://localhost:%s/api/docs", s.config.Port)
	log.Printf("❤️  Health Check: http://localhost:%s/health", s.config.Port)
	log.Printf("👤 Current User: Nayasora")
	log.Printf("🕐 Started at: 2025-06-02 15:00:58 UTC")
	log.Fatal(s.app.Listen(":" + s.config.Port))
}

func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":     message,
		"timestamp": "2025-06-02 15:00:58",
		"path":      c.Path(),
		"method":    c.Method(),
	})
}
