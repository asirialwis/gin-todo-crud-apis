package main

import (
	"gin-demo-api/db"
	"gin-demo-api/handlers"
	"gin-demo-api/middleware"

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

// --- JWT Security Definitions ---
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer <token>"

func main() {
	// 1. Initialize DB connection and run migrations
	db.ConnectDatabase()

	// 2. Initialize the Gin router
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Public Routes (Auth) ---
	router.POST("/register", handlers.RegisterUser) // Create a new user
	router.POST("/login", handlers.LoginUser)       // Log in and get JWT token

	// --- Protected Routes (Todos) ---
	todos := router.Group("/todos")
	todos.Use(middleware.RequireAuth) // Apply JWT middleware to all Todo routes
	{
		todos.POST("/", handlers.CreateTodo)
		todos.GET("/", handlers.FindTodos)
		todos.GET("/:id", handlers.FindTodo)
		todos.PATCH("/:id", handlers.UpdateTodo)
		todos.DELETE("/:id", handlers.DeleteTodo)
	}

	// --- User Routes (Protected) ---
	// Note: User CRUD endpoints now rely on the authenticated user and token
	users := router.Group("/users")
	users.Use(middleware.RequireAuth) // Apply JWT middleware to User routes
	{
		// Placeholder for user specific operations, usually these would be 'Me' endpoints
		// Example: users.GET("/me", handlers.GetUserProfile)
	}

	// 4. Start the server
	router.Run("localhost:8080")
}
