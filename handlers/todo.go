package handlers

import (
	"gin-demo-api/db"
	"gin-demo-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TodoRequestBody defines the request body for creating/updating a todo.
type TodoRequestBody struct {
	Item      string `json:"item" example:"Finish project report"`
	Completed *bool  `json:"completed" example:"false"` // Pointer to distinguish between false and not sent
}

// getAuthenticatedUserID retrieves the user ID from the Gin context, set by the RequireAuth middleware.
func getAuthenticatedUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication context missing"})
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return 0, false
	}
	return userID, true
}

// --- C R E A T E (POST /todos) ------------------------------------------------
// @Summary Create a new todo item
// @Description Creates a new todo item and links it to a user via user_id.
// @tags Todos
// @Accept  json
// @Produce  json
// @Param todo body models.Todo true "Todo item data (requires user_id)"
// @Success 201 {object} models.Todo
// @Failure 400 {object} map[string]interface{} "Invalid input format or invalid User ID"
// @Security ApiKeyAuth
// @Router /todos [post]
func CreateTodo(c *gin.Context) {

	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	var body TodoRequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default completed to false if not provided
	completed := false
	if body.Completed != nil {
		completed = *body.Completed
	}

	todo := models.Todo{Item: body.Item, Completed: completed, UserID: userID}
	result := db.DB.Create(&todo)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

// --- R E A D A L L (GET /todos) ---------------------------------------------
// @Summary Get all todo items
// @Description Retrieves a list of all todo items.
// @tags Todos
// @Produce  json
// @Success 200 {array} models.Todo
// @Security ApiKeyAuth
// @Router /todos [get]
func FindTodos(c *gin.Context) {

	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	var todos []models.Todo

	// Find all Todo records
	db.DB.Where("user_id = ?", userID).Find(&todos)

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
// @Security ApiKeyAuth
// @Router /todos/{id} [get]
func FindTodo(c *gin.Context) {

	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}

	id := c.Param("id")
	var todo models.Todo

	// Find todo by ID and check ownership
	result := db.DB.First(&todo, "id = ? AND user_id = ?", id, userID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or access denied"})
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
// @Security ApiKeyAuth
// @Router /todos/{id} [patch]
func UpdateTodo(c *gin.Context) {

	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var body TodoRequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo models.Todo
	// Find todo by ID and check ownership
	result := db.DB.First(&todo, "id = ? AND user_id = ?", id, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or access denied"})
		return
	}

	// Apply updates. Item check prevents updating to an empty string.
	if body.Item != "" {
		todo.Item = body.Item
	}
	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	db.DB.Save(&todo)
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
// @Security ApiKeyAuth
// @Router /todos/{id} [delete]
func DeleteTodo(c *gin.Context) {
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var todo models.Todo

	// Find todo by ID and check ownership before deleting
	result := db.DB.First(&todo, "id = ? AND user_id = ?", id, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or access denied"})
		return
	}

	// Use Soft Delete feature of GORM
	db.DB.Delete(&todo)
	c.Status(http.StatusNoContent)
}
