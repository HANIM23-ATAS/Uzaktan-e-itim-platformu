// Package database manages the SQLite connection and schema migrations.
package database

import (
	"log"

	"golearn/config"
	"golearn/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database handle.
var DB *gorm.DB

// Connect opens the SQLite file and runs AutoMigrate for all models.
func Connect(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Lesson{},
		&models.Quiz{},
		&models.Question{},
		&models.Progress{},
		&models.QuizResult{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("database connected and migrated")
}
