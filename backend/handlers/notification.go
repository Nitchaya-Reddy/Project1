package handlers

import (
	"net/http"
	"strconv"
	"time"
	"uf-marketplace/database"
	"uf-marketplace/models"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID")
	unreadOnly := c.DefaultQuery("unread", "false")

	query := database.DB.Where("user_id = ?", userID)

	if unreadOnly == "true" {
		query = query.Where("is_read = ?", false)
	}

	var notifications []models.Notification
	result := query.Order("created_at DESC").Find(&notifications)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	var count int64
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func MarkNotificationRead(c *gin.Context) {
	userID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	var notification models.Notification
	if result := database.DB.First(&notification, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	now := time.Now()
	database.DB.Model(&notification).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	now := time.Now()
	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func DeleteNotification(c *gin.Context) {
	userID := c.GetUint("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	var notification models.Notification
	if result := database.DB.First(&notification, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	database.DB.Delete(&notification)

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}
