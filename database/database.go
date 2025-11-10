package database

import (
	"log"

	"class-go-ai/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize database connection and auto-migrate models
func Connect() error {
	var err error
	
	// Open database connection
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	log.Println("Database connected and migrated successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
