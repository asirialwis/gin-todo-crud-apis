package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gin-demo-api/db"
	"gin-demo-api/models"
)

// --- Swagger Response Models ---

// ErrorResponse is the generic structure for failure responses.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid login data"`
}

// RegisterSuccessResponse defines the output structure for successful registration.
type RegisterSuccessResponse struct {
	Message  string `json:"message" example:"User registered successfully"`
	UserID   uint   `json:"user_id" example:"1"`
	Username string `json:"username" example:"user_alice"`
}

// LoginSuccessResponse defines the output structure for successful login.
type LoginSuccessResponse struct {
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserID uint   `json:"user_id" example:"1"`
}

// --- Request Body Definitions ---

// RegisterRequestBody defines the required fields for user registration.
type RegisterRequestBody struct {
	Username string `json:"username" binding:"required" example:"user_alice"`
	Email    string `json:"email" binding:"required,email" example:"alice@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePassword123"`
}

// LoginRequestBody defines the required fields for user login.
type LoginRequestBody struct {
	Email    string `json:"email" binding:"required,email" example:"alice@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePassword123"`
}

// --- Handlers ---

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a new user with a username, email, and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequestBody true "Registration details"
// @Success 201 {object} RegisterSuccessResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid input or email already in use"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var body RegisterRequestBody
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid registration data"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}

	// Create the User
	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email or Username already in use"})
		return
	}

	// Respond
	c.JSON(http.StatusCreated, RegisterSuccessResponse{
		Message:  "User registered successfully",
		UserID:   user.ID,
		Username: user.Username,
	})
}

// LoginUser godoc
// @Summary Login and get a JWT token
// @Description Authenticates a user with email and password and returns a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequestBody true "User credentials"
// @Success 200 {object} LoginSuccessResponse "Login successful, returns JWT token"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 400 {object} ErrorResponse "Invalid login data"
// @Failure 500 {object} ErrorResponse "Failed to generate token"
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var body LoginRequestBody
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid login data"})
		return
	}

	// Find the user by email
	var user models.User
	db.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Compare sent password with stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginSuccessResponse{Token: tokenString, UserID: user.ID})
}
