package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationNewMessage   NotificationType = "new_message"
	NotificationNewOffer     NotificationType = "new_offer"
	NotificationListingSold  NotificationType = "listing_sold"
	NotificationPriceDropped NotificationType = "price_dropped"
)

type Notification struct {
	ID        uint             `gorm:"primarykey" json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
	UserID    uint             `gorm:"not null" json:"user_id"`
	User      User             `gorm:"foreignKey:UserID" json:"-"`
	Type      NotificationType `gorm:"not null" json:"type"`
	Title     string           `gorm:"not null" json:"title"`
	Message   string           `json:"message"`
	Link      string           `json:"link"`
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
}
