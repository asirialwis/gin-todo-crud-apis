package handlers

import (
	"gin-demo-api/db"
	"gin-demo-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- C R E A T E (POST /todos) ------------------------------------------------
// @Summary Create a new todo item
// @Description Creates a new todo item and links it to a user via user_id.
// @tags Todos
// @Accept  json
// @Produce  json
// @Param todo body models.Todo true "Todo item data (requires user_id)"
// @Success 201 {object} models.Todo
// @Failure 400 {object} map[string]interface{} "Invalid input format or invalid User ID"
// @Router /todos [post]
func CreateTodo(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate UserID exists before creating todo
	var user models.User
	if err := db.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	// Save the new Todo record to the database
	db.DB.Create(&input)

	c.JSON(http.StatusCreated, input)
}

// --- R E A D A L L (GET /todos) ---------------------------------------------
// @Summary Get all todo items
// @Description Retrieves a list of all todo items.
// @tags Todos
// @Produce  json
// @Success 200 {array} models.Todo
// @Router /todos [get]
func FindTodos(c *gin.Context) {
	var todos []models.Todo

	// Find all Todo records
	db.DB.Find(&todos)

	c.JSON(http.StatusOK, todos)
}

// --- R E A D O N E (GET /todos/:id) -----------------------------------------
// @Summary Get todo item by ID
// @Description Retrieves a single todo item by its ID.
// @tags Todos
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Router /todos/{id} [get]
func FindTodo(c *gin.Context) {
	var todo models.Todo

	// Find record by ID (from URL parameter)
	if err := db.DB.Where("id = ?", c.Param("id")).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// --- U P D A T E (PATCH /todos/:id) -----------------------------------------
// @Summary Update a todo item
// @Description Updates the item and/or completed status for a specific todo.
// @tags Todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo data (item and/or completed status)"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]interface{} "Invalid input format"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Router /todos/{id} [patch]
func UpdateTodo(c *gin.Context) {
	var todo models.Todo
	// Check if todo exists
	if err := db.DB.Where("id = ?", c.Param("id")).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var input models.Todo
	// Validate input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the record with the new input data
	db.DB.Model(&todo).Updates(input)

	c.JSON(http.StatusOK, todo)
}

// --- D E L E T E (DELETE /todos/:id) ----------------------------------------
// @Summary Delete a todo item
// @Description Soft-deletes a todo item by ID.
// @tags Todos
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} map[string]interface{} "Deletion successful"
// @Failure 404 {object} map[string]interface{} "Todo not found"
// @Router /todos/{id} [delete]
func DeleteTodo(c *gin.Context) {
	var todo models.Todo
	// Check if todo exists
	if err := db.DB.Where("id = ?", c.Param("id")).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// Soft delete the record
	db.DB.Delete(&todo)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
