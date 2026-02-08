package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	ChatID     uint           `gorm:"not null" json:"chat_id"`
	Chat       Chat           `gorm:"foreignKey:ChatID" json:"-"`
	SenderID   uint           `gorm:"not null" json:"sender_id"`
	Sender     User           `gorm:"foreignKey:SenderID" json:"sender"`
	Content    string         `gorm:"not null" json:"content"`
	IsRead     bool           `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time     `json:"read_at,omitempty"`
}

type Chat struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	ListingID  uint           `gorm:"not null" json:"listing_id"`
	Listing    Listing        `gorm:"foreignKey:ListingID" json:"listing"`
	BuyerID    uint           `gorm:"not null" json:"buyer_id"`
	Buyer      User           `gorm:"foreignKey:BuyerID" json:"buyer"`
	SellerID   uint           `gorm:"not null" json:"seller_id"`
	Seller     User           `gorm:"foreignKey:SellerID" json:"seller"`
	Messages   []Message      `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
	LastMessage *Message      `gorm:"-" json:"last_message,omitempty"`
}

type ChatResponse struct {
	ID          uint         `json:"id"`
	ListingID   uint         `json:"listing_id"`
	Listing     Listing      `json:"listing"`
	BuyerID     uint         `json:"buyer_id"`
	Buyer       UserResponse `json:"buyer"`
	SellerID    uint         `json:"seller_id"`
	Seller      UserResponse `json:"seller"`
	LastMessage *Message     `json:"last_message,omitempty"`
	UnreadCount int          `json:"unread_count"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
