package database

import (
	"log"
	"uf-marketplace/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("marketplace.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Listing{},
		&models.ListingImage{},
		&models.Chat{},
		&models.Message{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed categories if they don't exist
	seedCategories()

	log.Println("Database initialized successfully")
}

func seedCategories() {
	categories := []models.Category{
		{Name: "Textbooks", Description: "Academic textbooks and study materials", Icon: "book"},
		{Name: "Electronics", Description: "Phones, laptops, tablets, and accessories", Icon: "devices"},
		{Name: "Furniture", Description: "Dorm and apartment furniture", Icon: "chair"},
		{Name: "Clothing", Description: "Clothes, shoes, and accessories", Icon: "checkroom"},
		{Name: "Sports", Description: "Sports equipment and gear", Icon: "sports_soccer"},
		{Name: "Tickets", Description: "Event and game tickets", Icon: "confirmation_number"},
		{Name: "Transportation", Description: "Bikes, scooters, and car accessories", Icon: "directions_bike"},
		{Name: "Services", Description: "Tutoring, moving help, etc.", Icon: "handyman"},
		{Name: "Housing", Description: "Sublease and roommate listings", Icon: "home"},
		{Name: "Other", Description: "Everything else", Icon: "category"},
	}

	for _, category := range categories {
		DB.FirstOrCreate(&category, models.Category{Name: category.Name})
	}
}

func GetDB() *gorm.DB {
	return DB
}
