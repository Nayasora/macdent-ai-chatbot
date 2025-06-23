package databases

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"macdent-ai-chatbot/v2/configs"
	"macdent-ai-chatbot/v2/models"
	"macdent-ai-chatbot/v2/utils"
)

type PostgresDatabase struct {
	DB *gorm.DB
}

func NewPostgres(config *configs.PostgresConfig) *PostgresDatabase {
	logger := utils.NewLogger("postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		config.Host, config.User, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("postgres клиент: %v", err)
	}

	err = models.InitMigration(db)

	if err != nil {
		logger.Errorf("postgres миграция: %v", err)
	}

	return &PostgresDatabase{DB: db}
}
