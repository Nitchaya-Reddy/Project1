package models

import (
	"time"

	"gorm.io/gorm"
)

type ListingStatus string

const (
	StatusActive   ListingStatus = "active"
	StatusSold     ListingStatus = "sold"
	StatusInactive ListingStatus = "inactive"
)

type Listing struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	CategoryID  uint           `json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	SellerID    uint           `gorm:"not null" json:"seller_id"`
	Seller      User           `gorm:"foreignKey:SellerID" json:"seller"`
	Images      []ListingImage `gorm:"foreignKey:ListingID" json:"images"`
	Status      ListingStatus  `gorm:"default:'active'" json:"status"`
	Condition   string         `json:"condition"` // new, like_new, good, fair, poor
	Location    string         `json:"location"`
	Views       int            `gorm:"default:0" json:"views"`
}

type ListingImage struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ListingID uint           `gorm:"not null" json:"listing_id"`
	ImageURL  string         `gorm:"not null" json:"image_url"`
	IsPrimary bool           `gorm:"default:false" json:"is_primary"`
}

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
}
