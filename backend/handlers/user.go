package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"uf-marketplace/database"
	"uf-marketplace/models"
	"uf-marketplace/utils"

	"github.com/gin-gonic/gin"
)

type UpdateUserInput struct {
	Name         string `json:"name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func UpdateUser(c *gin.Context) {
	userID := c.GetUint("userID")

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided
	if input.Name != "" {
		// Split name into first and last
		parts := strings.SplitN(input.Name, " ", 2)
		user.FirstName = parts[0]
		if len(parts) > 1 {
			user.LastName = parts[1]
		}
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	if input.Bio != "" {
		user.Bio = input.Bio
	}
	if input.ProfileImage != "" {
		user.ProfileImage = input.ProfileImage
	}

	if result := database.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func GetUserListings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var listings []models.Listing
	result := database.DB.
		Preload("Images").
		Preload("Category").
		Preload("Seller").
		Where("seller_id = ?", id).
		Order("created_at DESC").
		Find(&listings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching listings"})
		return
	}

	c.JSON(http.StatusOK, listings)
}

func GetMyListings(c *gin.Context) {
	userID := c.GetUint("userID")

	status := c.DefaultQuery("status", "")

	query := database.DB.
		Preload("Images").
		Preload("Category").
		Preload("Seller").
		Where("seller_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var listings []models.Listing
	result := query.Order("created_at DESC").Find(&listings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching listings"})
		return
	}

	c.JSON(http.StatusOK, listings)
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	if !utils.CheckPassword(input.CurrentPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password"})
		return
	}

	user.Password = hashedPassword
	if result := database.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
