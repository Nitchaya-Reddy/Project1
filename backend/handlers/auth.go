package handlers

import (
	"net/http"
	"strings"
	"uf-marketplace/database"
	"uf-marketplace/models"
	"uf-marketplace/utils"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Parse validation errors for better messages
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "Email") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a valid email address"})
			return
		}
		if strings.Contains(errorMsg, "Password") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
			return
		}
		if strings.Contains(errorMsg, "FirstName") || strings.Contains(errorMsg, "LastName") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "First name and last name are required"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill in all required fields"})
		return
	}

	// Check if email is UF email (optional - for UF students only)
	if !strings.HasSuffix(strings.ToLower(input.Email), "@ufl.edu") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must use a valid UF email (@ufl.edu)"})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if result := database.DB.Where("email = ?", input.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "An account with this email already exists. Please login or use a different email."})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
		return
	}

	// Create user
	user := models.User{
		Email:     strings.ToLower(input.Email),
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	if result := database.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "Email") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a valid email address"})
			return
		}
		if strings.Contains(errorMsg, "Password") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter email and password"})
		return
	}

	var user models.User
	if result := database.DB.Where("email = ?", strings.ToLower(input.Email)).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password. Please try again."})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password. Please try again."})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func GetMe(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
