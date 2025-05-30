package config

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"macdent-ai-chatbot/logger"
	"os"
	"strconv"
)

type Env struct {
	logger   *log.Logger
	filename string
}

// NewEnv создает новый экземпляр сервиса и загружает .env файл
func NewEnv(filename string) *Env {
	customLogger := logger.New("env")

	err := godotenv.Load(filename)

	if err != nil {
		customLogger.Fatalf("загрузка файла: %s", err)
	}

	customLogger.Infof("файл загружен: %s", filename)

	return &Env{
		logger:   customLogger,
		filename: filename,
	}
}

// MustString возвращает строковое значение или завершает работу с ошибкой
func (e *Env) MustString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		e.logger.Fatalf("%s: обязательная переменная - %s", e.filename, key)
	}
	return value
}

// MustInt возвращает числовое значение или завершает работу с ошибкой
func (e *Env) MustInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		e.logger.Fatalf("%s: обязательная переменная - %s", e.filename, key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		e.logger.Fatalf("%s: переменная %s должна быть числом: %v", e.filename, key, err)
	}

	return intValue
}

// MustBool возвращает булево значение или завершает работу с ошибкой
func (e *Env) MustBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		e.logger.Fatalf("%s: обязательная переменная - %s", e.filename, key)
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		e.logger.Fatalf("%s: переменная %s должна быть булевым: %v", e.filename, key, err)
	}

	return boolValue
}
