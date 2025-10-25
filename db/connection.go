package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gin-demo-api/models"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection and runs migrations
func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	// AutoMigrate creates the table based on the Todo struct
	err = database.AutoMigrate(&models.Todo{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to run migrations!")
	}

	DB = database
}
