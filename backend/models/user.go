package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"`
	FirstName    string         `gorm:"not null" json:"first_name"`
	LastName     string         `gorm:"not null" json:"last_name"`
	ProfileImage string         `json:"profile_image"`
	Phone        string         `json:"phone"`
	Bio          string         `json:"bio"`
	IsAdmin      bool           `gorm:"default:false" json:"is_admin"`
	Listings     []Listing      `gorm:"foreignKey:SellerID" json:"listings,omitempty"`
	Messages     []Message      `gorm:"foreignKey:SenderID" json:"messages,omitempty"`
}

type UserResponse struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	ProfileImage string    `json:"profile_image"`
	Phone        string    `json:"phone"`
	Bio          string    `json:"bio"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		ProfileImage: u.ProfileImage,
		Phone:        u.Phone,
		Bio:          u.Bio,
		IsAdmin:      u.IsAdmin,
		CreatedAt:    u.CreatedAt,
	}
}
