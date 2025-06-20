package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"macdent-ai-chatbot/v1/internal/models"
	"time"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(host, user, password, dbname string, port int) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Agent{}, &models.Dialog{}, &models.KnowledgeFile{})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) CreateAgent(agent *models.Agent) error {
	return d.DB.Create(agent).Error
}

func (d *Database) GetAgentByID(id string) (*models.Agent, error) {
	var agent models.Agent
	err := d.DB.First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (d *Database) UpdateAgent(agent *models.Agent) error {
	return d.DB.Save(agent).Error
}

func (d *Database) DeleteAgent(id string) error {
	return d.DB.Delete(&models.Agent{}, "id = ?", id).Error
}

func (d *Database) ListAgents(limit, offset int) ([]models.Agent, error) {
	var agents []models.Agent
	err := d.DB.Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&agents).Error
	return agents, err
}

func (d *Database) CreateDialog(dialog *models.Dialog) error {
	fmt.Println(dialog)
	return d.DB.Create(dialog).Error
}

func (d *Database) GetDialogHistory(agentID, userID string, limit int) ([]models.Dialog, error) {
	var dialogs []models.Dialog
	err := d.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).
		Order("created_at ASC").
		Limit(limit).
		Find(&dialogs).Error
	return dialogs, err
}

func (d *Database) GetDialogsByAgent(agentID string, limit, offset int) ([]models.Dialog, error) {
	var dialogs []models.Dialog
	err := d.DB.Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&dialogs).Error
	return dialogs, err
}

func (d *Database) DeleteDialogHistory(agentID, userID string) error {
	return d.DB.Where("agent_id = ? AND user_id = ?", agentID, userID).
		Delete(&models.Dialog{}).Error
}

func (d *Database) CountDialogsByAgent(agentID string) (int64, error) {
	var count int64
	err := d.DB.Model(&models.Dialog{}).Where("agent_id = ?", agentID).Count(&count).Error
	return count, err
}

func (d *Database) GetLastDialogTime(agentID string) (time.Time, error) {
	var dialog models.Dialog
	err := d.DB.Where("agent_id = ?", agentID).
		Order("created_at DESC").
		First(&dialog).Error
	if err != nil {
		return time.Time{}, err
	}
	return dialog.CreatedAt, nil
}

func (d *Database) CreateKnowledgeFile(file *models.KnowledgeFile) error {
	return d.DB.Create(file).Error
}

func (d *Database) GetKnowledgeFiles(agentID string) ([]models.KnowledgeFile, error) {
	var files []models.KnowledgeFile
	err := d.DB.Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

func (d *Database) GetKnowledgeFileByID(id uint) (*models.KnowledgeFile, error) {
	var file models.KnowledgeFile
	err := d.DB.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (d *Database) UpdateKnowledgeFile(file *models.KnowledgeFile) error {
	return d.DB.Save(file).Error
}

func (d *Database) DeleteKnowledgeFile(id uint) error {
	return d.DB.Delete(&models.KnowledgeFile{}, id).Error
}

func (d *Database) DeleteKnowledgeFilesByAgent(agentID string) error {
	return d.DB.Where("agent_id = ?", agentID).Delete(&models.KnowledgeFile{}).Error
}

func (d *Database) UpdateKnowledgeFileStatus(id uint, status string, errorMessage string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "completed" {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if errorMessage != "" {
		updates["error_message"] = errorMessage
	}

	return d.DB.Model(&models.KnowledgeFile{}).Where("id = ?", id).Updates(updates).Error
}

func (d *Database) HealthCheck() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
