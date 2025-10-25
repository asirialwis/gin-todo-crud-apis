package main

import (
	"gin-demo-api/db"
	"gin-demo-api/handlers"

	"github.com/gin-gonic/gin"

	// imports for Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// 🚨 Import the docs package (must be manually created by 'swag init')
	_ "gin-demo-api/docs"
)

// 🚨 Add top-level annotations for the API metadata
// @title Gin CRUD API
// @version 1.0
// @description This is a sample server for a User/Todo management API.
// @host localhost:8080
// @BasePath /

func main() {
	// 1. Initialize DB connection and run migrations
	db.ConnectDatabase()

	// 2. Initialize the Gin router
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- USER ROUTES ---
	router.POST("/users", handlers.CreateUser)       // C: Create User
	router.GET("/users", handlers.FindUsers)         // R: Read All Users (with Todos)
	router.GET("/users/:id", handlers.FindUser)      // R: Read One User (with Todos)
	router.PATCH("/users/:id", handlers.UpdateUser)  // U: Update User
	router.DELETE("/users/:id", handlers.DeleteUser) // D: Delete User

	// 3. Define RESTful API routes (CRUD)
	router.POST("/todos", handlers.CreateTodo)       // C: Create
	router.GET("/todos", handlers.FindTodos)         // R: Read All
	router.GET("/todos/:id", handlers.FindTodo)      // R: Read One
	router.PATCH("/todos/:id", handlers.UpdateTodo)  // U: Update
	router.DELETE("/todos/:id", handlers.DeleteTodo) // D: Delete

	// 4. Start the server
	router.Run("localhost:8080")
}
