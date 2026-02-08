package handlers

import (
	"net/http"
	"strconv"
	"time"
	"uf-marketplace/database"
	"uf-marketplace/models"

	"github.com/gin-gonic/gin"
)

type CreateChatInput struct {
	ListingID uint   `json:"listing_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type SendMessageInput struct {
	Content string `json:"content" binding:"required"`
}

func GetChats(c *gin.Context) {
	userID := c.GetUint("userID")

	var chats []models.Chat
	result := database.DB.
		Preload("Listing").
		Preload("Listing.Images").
		Preload("Buyer").
		Preload("Seller").
		Where("buyer_id = ? OR seller_id = ?", userID, userID).
		Order("updated_at DESC").
		Find(&chats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching chats"})
		return
	}

	// Get last message and unread count for each chat
	chatResponses := make([]models.ChatResponse, 0)
	for _, chat := range chats {
		var lastMessage models.Message
		database.DB.Where("chat_id = ?", chat.ID).Order("created_at DESC").First(&lastMessage)

		var unreadCount int64
		database.DB.Model(&models.Message{}).
			Where("chat_id = ? AND sender_id != ? AND is_read = ?", chat.ID, userID, false).
			Count(&unreadCount)

		response := models.ChatResponse{
			ID:          chat.ID,
			ListingID:   chat.ListingID,
			Listing:     chat.Listing,
			BuyerID:     chat.BuyerID,
			Buyer:       chat.Buyer.ToResponse(),
			SellerID:    chat.SellerID,
			Seller:      chat.Seller.ToResponse(),
			UnreadCount: int(unreadCount),
			CreatedAt:   chat.CreatedAt,
			UpdatedAt:   chat.UpdatedAt,
		}

		if lastMessage.ID > 0 {
			response.LastMessage = &lastMessage
		}

		chatResponses = append(chatResponses, response)
	}

	c.JSON(http.StatusOK, chatResponses)
}

func CreateChat(c *gin.Context) {
	userID := c.GetUint("userID")

	var input CreateChatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get listing
	var listing models.Listing
	if result := database.DB.First(&listing, input.ListingID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	// Can't message your own listing
	if listing.SellerID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot message your own listing"})
		return
	}

	// Check if chat already exists
	var existingChat models.Chat
	result := database.DB.
		Where("listing_id = ? AND buyer_id = ?", input.ListingID, userID).
		First(&existingChat)

	if result.Error == nil {
		// Chat exists, just add message
		message := models.Message{
			ChatID:   existingChat.ID,
			SenderID: userID,
			Content:  input.Message,
		}
		database.DB.Create(&message)
		database.DB.Model(&existingChat).Update("updated_at", time.Now())

		c.JSON(http.StatusOK, gin.H{
			"chat_id": existingChat.ID,
			"message": message,
		})
		return
	}

	// Create new chat
	chat := models.Chat{
		ListingID: input.ListingID,
		BuyerID:   userID,
		SellerID:  listing.SellerID,
	}

	if result := database.DB.Create(&chat); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chat"})
		return
	}

	// Create first message
	message := models.Message{
		ChatID:   chat.ID,
		SenderID: userID,
		Content:  input.Message,
	}
	database.DB.Create(&message)

	// Create notification for seller
	notification := models.Notification{
		UserID:  listing.SellerID,
		Type:    models.NotificationNewMessage,
		Title:   "New Message",
		Message: "You have a new message about your listing: " + listing.Title,
		Link:    "/chat/" + strconv.Itoa(int(chat.ID)),
	}
	database.DB.Create(&notification)

	c.JSON(http.StatusCreated, gin.H{
		"chat_id": chat.ID,
		"message": message,
	})
}

func GetChat(c *gin.Context) {
	userID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var chat models.Chat
	result := database.DB.
		Preload("Listing").
		Preload("Listing.Images").
		Preload("Buyer").
		Preload("Seller").
		First(&chat, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Check if user is part of the chat
	if chat.BuyerID != userID && chat.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func GetChatMessages(c *gin.Context) {
	userID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var chat models.Chat
	if result := database.DB.First(&chat, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Check if user is part of the chat
	if chat.BuyerID != userID && chat.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this chat"})
		return
	}

	var messages []models.Message
	database.DB.
		Preload("Sender").
		Where("chat_id = ?", id).
		Order("created_at ASC").
		Find(&messages)

	// Mark messages as read
	now := time.Now()
	database.DB.Model(&models.Message{}).
		Where("chat_id = ? AND sender_id != ? AND is_read = ?", id, userID, false).
		Updates(map[string]interface{}{"is_read": true, "read_at": now})

	c.JSON(http.StatusOK, messages)
}

func SendMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	var input SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var chat models.Chat
	if result := database.DB.Preload("Listing").First(&chat, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Check if user is part of the chat
	if chat.BuyerID != userID && chat.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to send messages in this chat"})
		return
	}

	message := models.Message{
		ChatID:   uint(id),
		SenderID: userID,
		Content:  input.Content,
	}

	if result := database.DB.Create(&message); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message"})
		return
	}

	// Update chat timestamp
	database.DB.Model(&chat).Update("updated_at", time.Now())

	// Create notification for recipient
	var recipientID uint
	if chat.BuyerID == userID {
		recipientID = chat.SellerID
	} else {
		recipientID = chat.BuyerID
	}

	notification := models.Notification{
		UserID:  recipientID,
		Type:    models.NotificationNewMessage,
		Title:   "New Message",
		Message: "You have a new message about: " + chat.Listing.Title,
		Link:    "/chat/" + strconv.Itoa(int(chat.ID)),
	}
	database.DB.Create(&notification)

	// Reload with sender
	database.DB.Preload("Sender").First(&message, message.ID)

	c.JSON(http.StatusCreated, message)
}
