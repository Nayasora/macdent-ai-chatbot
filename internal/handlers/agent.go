package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"macdent-ai-chatbot/internal/services"
)

type AgentHandler struct {
	agentService     *services.AgentService
	dialogService    *services.DialogService
	knowledgeService *services.KnowledgeService
}

func NewAgentHandler(
	agentService *services.AgentService,
	dialogService *services.DialogService,
	knowledgeService *services.KnowledgeService,
) *AgentHandler {
	return &AgentHandler{
		agentService:     agentService,
		dialogService:    dialogService,
		knowledgeService: knowledgeService,
	}
}

func (h *AgentHandler) CreateAgent(c fiber.Ctx) error {
	var req services.CreateAgentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	agent, err := h.agentService.CreateAgent(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create agent",
			"details": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Agent created successfully",
		"agent":   agent,
	})
}

func (h *AgentHandler) GetAgent(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	agent, err := h.agentService.GetAgent(agentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Agent not found",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"agent": agent,
	})
}

func (h *AgentHandler) UpdateAgent(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	var req services.UpdateAgentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	agent, err := h.agentService.UpdateAgent(agentID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update agent",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Agent updated successfully",
		"agent":   agent,
	})
}

func (h *AgentHandler) DeleteAgent(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	err = h.agentService.DeleteAgent(agentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete agent",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Agent deleted successfully",
	})
}

func (h *AgentHandler) ListAgents(c fiber.Ctx) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	agents, err := h.agentService.ListAgents(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to list agents",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"agents": agents,
		"pagination": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *AgentHandler) GetAgentStats(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	stats, err := h.agentService.GetAgentStats(agentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get agent stats",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
	})
}

func (h *AgentHandler) UploadKnowledge(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	// Проверяем существование агента
	_, err = h.agentService.GetAgent(agentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Agent not found",
			"details": err.Error(),
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}

	fmt.Println("Received multipart form:", form)

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No files provided",
		})
	}

	uploadedFiles := make([]map[string]interface{}, 0, len(files))
	errors := make([]string, 0)

	for _, file := range files {
		result, err := h.processUploadedFile(agentID, file)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		uploadedFiles = append(uploadedFiles, result)
	}

	response := fiber.Map{
		"message":        "File upload completed",
		"uploaded_files": uploadedFiles,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	statusCode := 201
	if len(uploadedFiles) == 0 {
		statusCode = 400
		response["message"] = "No files were successfully uploaded"
	} else if len(errors) > 0 {
		statusCode = 207
		response["message"] = "Some files failed to upload"
	}

	return c.Status(statusCode).JSON(response)
}

func (h *AgentHandler) GetKnowledge(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	knowledgeFiles, err := h.knowledgeService.GetKnowledgeFiles(agentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to retrieve knowledge files",
			"details": err.Error(),
		})
	}

	stats, err := h.knowledgeService.GetKnowledgeStats(agentID)
	if err != nil {
		stats = map[string]interface{}{}
	}

	return c.JSON(fiber.Map{
		"files": knowledgeFiles,
		"stats": stats,
	})
}

func (h *AgentHandler) DeleteKnowledge(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	fileID := c.Query("file_id")
	if fileID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "file_id query parameter is required",
		})
	}

	uintFileID, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		fmt.Println("Error parsing file_id:", err)
		os.Exit(1)
	}

	err = h.knowledgeService.DeleteKnowledgeFile(agentID, uint(uintFileID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete knowledge file",
			"details": err.Error(),
		})
	}

	// Удаляем файл из файловой системы
	filePath := filepath.Join("uploads", agentID.String(), fileID)
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		os.Exit(1)
	}

	return c.JSON(fiber.Map{
		"message": "Knowledge file deleted successfully",
	})
}

func (h *AgentHandler) ProcessDialog(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	var req services.DialogRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	response, err := h.dialogService.ProcessDialog(agentID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to process dialog",
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *AgentHandler) GetDialogHistory(c fiber.Ctx) error {
	idParam := c.Params("id")
	agentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid agent ID format",
		})
	}

	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "user_id query parameter is required",
		})
	}

	limitStr := c.Query("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	dialogs, err := h.dialogService.GetDialogHistory(agentID, userID, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get dialog history",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"dialogs": dialogs,
		"user_id": userID,
		"limit":   limit,
	})
}

func (h *AgentHandler) processUploadedFile(agentID uuid.UUID, file *multipart.FileHeader) (map[string]interface{}, error) {
	// Создаем директорию для хранения файлов
	uploadDir := filepath.Join("uploads", agentID.String())
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return nil, err
	}

	// Генерируем уникальное имя файла
	filename := generateUniqueFilename(file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Сохраняем файл
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
			os.Exit(1)
		}
	}(src)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
			os.Exit(1)
		}
	}(dst)

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	// Читаем содержимое файла
	content, err := h.readFileContent(filePath, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	// Добавляем в базу знаний
	err = h.knowledgeService.AddDocument(
		agentID,
		filename,
		content,
		file.Size,
		file.Header.Get("Content-Type"),
	)
	if err != nil {
		// Удаляем файл при ошибке
		err := os.Remove(filePath)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return map[string]interface{}{
		"filename":      filename,
		"original_name": file.Filename,
		"size":          file.Size,
		"type":          file.Header.Get("Content-Type"),
		"path":          filePath,
	}, nil
}

func (h *AgentHandler) readFileContent(filePath, contentType string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Простая обработка текстовых файлов
	// Можно расширить для поддержки PDF, DOCX и других форматов
	if strings.HasPrefix(contentType, "text/") || contentType == "application/json" {
		return string(content), nil
	}

	// Для остальных типов возвращаем как есть
	return string(content), nil
}

func generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	base := strings.TrimSuffix(originalName, ext)
	timestamp := time.Now().Format("20060102_150405")
	return base + "_" + timestamp + ext
}
