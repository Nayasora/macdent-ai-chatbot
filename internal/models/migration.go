package models

import (
	"gorm.io/gorm"
)

func InitMigration(db *gorm.DB) error {
	return db.AutoMigrate(
		&Agent{},
		&Permission{},
		&KnowledgePrompt{},
		&KnowledgeFile{},
		&Dialog{},
	)
}
