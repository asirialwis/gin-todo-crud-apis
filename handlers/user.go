package handlers

import (
	"gin-demo-api/db"
	"gin-demo-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- C R E A T E (POST /users) ------------------------------------------------
// @Summary Create a new user
// @Description Creates a new user with a unique username and email.
// @tags Users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User data (only username and email are required)"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]interface{} "Invalid input format or duplicate entry"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the new User record to the database
	db.DB.Create(&input)

	c.JSON(http.StatusCreated, input)
}

// --- R E A D A L L (GET /users) ---------------------------------------------
// @Summary Get all users
// @Description Retrieves a list of all users, preloading their associated todos.
// @tags Users
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get] // <-- CORRECT: /users [get] for ALL users
func FindUsers(c *gin.Context) {
	var users []models.User

	// Preload the Todos relationship when retrieving users
	db.DB.Preload("Todos").Find(&users)

	c.JSON(http.StatusOK, users)
}

// --- R E A D O N E (GET /users/:id) -----------------------------------------
// @Summary Get user by ID
// @Description Retrieves a single user by their ID.
// @tags Users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get] // <-- CORRECT: /users/{id} [get] for ONE user
func FindUser(c *gin.Context) {
	var user models.User

	// Find record by ID (from URL parameter), Preload Todos
	if err := db.DB.Preload("Todos").Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// --- U P D A T E (PATCH /users/:id) -----------------------------------------
// @Summary Update a user
// @Description Updates the username and/or email for a specific user.
// @tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.User true "User data (only username/email are updated)"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{} "Invalid input format"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [patch] // <-- CORRECT: /users/{id} [patch] for UPDATE
func UpdateUser(c *gin.Context) {
	var user models.User
	// Check if user exists
	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	// Validate input JSON, ignoring the Todos field
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the record with the new input data
	db.DB.Model(&user).Updates(models.User{Username: input.Username, Email: input.Email})

	c.JSON(http.StatusOK, user)
}

// --- D E L E T E (DELETE /users/:id) ----------------------------------------
// @Summary Delete a user
// @Description Soft-deletes a user by ID.
// @tags Users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Deletion successful" // <-- FIXED gin.H here
// @Failure 404 {object} map[string]interface{} "User not found" // <-- FIXED gin.H here
// @Router /users/{id} [delete] // <-- CORRECT: /users/{id} [delete] for DELETE
func DeleteUser(c *gin.Context) {
	var user models.User
	// Check if user exists
	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// WARNING: In a real app, you must decide how to handle the dependent todos (e.g., delete them too, or set UserID to null)
	// For this demo, GORM will typically handle the soft delete on the User record.
	db.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
