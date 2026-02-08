package main

import (
	"log"
	"os"
	"strings"
	"uf-marketplace/database"
	"uf-marketplace/handlers"
	"uf-marketplace/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	if allowedOrigins != "" {
		config.AllowOrigins = strings.Split(allowedOrigins, ",")
	} else {
		config.AllowOrigins = []string{"http://localhost:4200", "http://localhost:3000"}
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
		}

		// Categories (public)
		api.GET("/categories", handlers.GetCategories)

		// Listings routes
		listings := api.Group("/listings")
		{
			listings.GET("", middleware.OptionalAuthMiddleware(), handlers.GetListings)
			listings.GET("/:id", middleware.OptionalAuthMiddleware(), handlers.GetListing)
			listings.POST("", middleware.AuthMiddleware(), handlers.CreateListing)
			listings.PUT("/:id", middleware.AuthMiddleware(), handlers.UpdateListing)
			listings.DELETE("/:id", middleware.AuthMiddleware(), handlers.DeleteListing)
		}

		// Upload route
		api.POST("/upload", middleware.AuthMiddleware(), handlers.UploadImage)

		// User routes
		users := api.Group("/users")
		{
			users.GET("/:id", handlers.GetUser)
			users.GET("/:id/listings", handlers.GetUserListings)
			users.PUT("/me", middleware.AuthMiddleware(), handlers.UpdateUser)
			users.PUT("/me/password", middleware.AuthMiddleware(), handlers.ChangePassword)
			users.GET("/me/listings", middleware.AuthMiddleware(), handlers.GetMyListings)
		}

		// Chat routes
		chats := api.Group("/chats")
		chats.Use(middleware.AuthMiddleware())
		{
			chats.GET("", handlers.GetChats)
			chats.POST("", handlers.CreateChat)
			chats.GET("/:id", handlers.GetChat)
			chats.GET("/:id/messages", handlers.GetChatMessages)
			chats.POST("/:id/messages", handlers.SendMessage)
		}

		// Notification routes
		notifications := api.Group("/notifications")
		notifications.Use(middleware.AuthMiddleware())
		{
			notifications.GET("", handlers.GetNotifications)
			notifications.GET("/unread-count", handlers.GetUnreadCount)
			notifications.PUT("/:id/read", handlers.MarkNotificationRead)
			notifications.PUT("/read-all", handlers.MarkAllNotificationsRead)
			notifications.DELETE("/:id", handlers.DeleteNotification)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
