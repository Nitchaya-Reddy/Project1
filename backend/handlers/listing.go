package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"uf-marketplace/database"
	"uf-marketplace/models"

	"github.com/gin-gonic/gin"
)

type CreateListingInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gte=0"`
	CategoryID  uint     `json:"category_id" binding:"required"`
	Condition   string   `json:"condition"`
	Location    string   `json:"location"`
	Images      []string `json:"images"`
}

type UpdateListingInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  uint     `json:"category_id"`
	Condition   string   `json:"condition"`
	Location    string   `json:"location"`
	Status      string   `json:"status"`
	Images      []string `json:"images"`
}

func CreateListing(c *gin.Context) {
	userID := c.GetUint("userID")

	var input CreateListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify category exists
	var category models.Category
	if result := database.DB.First(&category, input.CategoryID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}

	listing := models.Listing{
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		SellerID:    userID,
		Condition:   input.Condition,
		Location:    input.Location,
		Status:      models.StatusActive,
	}

	if result := database.DB.Create(&listing); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating listing"})
		return
	}

	// Add images
	for i, imageURL := range input.Images {
		image := models.ListingImage{
			ListingID: listing.ID,
			ImageURL:  imageURL,
			IsPrimary: i == 0,
		}
		database.DB.Create(&image)
	}

	// Reload with associations
	database.DB.Preload("Images").Preload("Category").Preload("Seller").First(&listing, listing.ID)

	c.JSON(http.StatusCreated, listing)
}

func GetListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	categoryID := c.Query("category_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	condition := c.Query("condition")
	sortBy := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "desc")

	offset := (page - 1) * limit

	query := database.DB.Model(&models.Listing{}).Where("status = ?", models.StatusActive)

	// Apply filters
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if condition != "" {
		query = query.Where("condition = ?", condition)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Apply sorting and pagination
	var listings []models.Listing
	result := query.
		Preload("Images").
		Preload("Category").
		Preload("Seller").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Offset(offset).
		Limit(limit).
		Find(&listings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching listings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listings": listings,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + int64(limit) - 1) / int64(limit),
	})
}

func GetListing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var listing models.Listing
	result := database.DB.
		Preload("Images").
		Preload("Category").
		Preload("Seller").
		First(&listing, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	// Increment view count
	database.DB.Model(&listing).Update("views", listing.Views+1)

	c.JSON(http.StatusOK, listing)
}

func UpdateListing(c *gin.Context) {
	userID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var listing models.Listing
	if result := database.DB.First(&listing, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	// Check ownership
	if listing.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this listing"})
		return
	}

	var input UpdateListingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if input.Title != "" {
		listing.Title = input.Title
	}
	if input.Description != "" {
		listing.Description = input.Description
	}
	if input.Price > 0 {
		listing.Price = input.Price
	}
	if input.CategoryID > 0 {
		listing.CategoryID = input.CategoryID
	}
	if input.Condition != "" {
		listing.Condition = input.Condition
	}
	if input.Location != "" {
		listing.Location = input.Location
	}
	if input.Status != "" {
		listing.Status = models.ListingStatus(input.Status)
	}

	if result := database.DB.Save(&listing); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating listing"})
		return
	}

	// Update images if provided
	if len(input.Images) > 0 {
		database.DB.Where("listing_id = ?", listing.ID).Delete(&models.ListingImage{})
		for i, imageURL := range input.Images {
			image := models.ListingImage{
				ListingID: listing.ID,
				ImageURL:  imageURL,
				IsPrimary: i == 0,
			}
			database.DB.Create(&image)
		}
	}

	database.DB.Preload("Images").Preload("Category").Preload("Seller").First(&listing, listing.ID)

	c.JSON(http.StatusOK, listing)
}

func DeleteListing(c *gin.Context) {
	userID := c.GetUint("userID")
	isAdmin := c.GetBool("isAdmin")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var listing models.Listing
	if result := database.DB.First(&listing, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	// Check ownership or admin
	if listing.SellerID != userID && !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this listing"})
		return
	}

	// Delete images first
	database.DB.Where("listing_id = ?", listing.ID).Delete(&models.ListingImage{})

	// Delete listing
	if result := database.DB.Delete(&listing); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image provided"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      "/uploads/" + filename,
		"filename": filename,
	})
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if result := database.DB.Find(&categories); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
