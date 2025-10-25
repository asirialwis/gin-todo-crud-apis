package models

import (
	"time"

	"gorm.io/gorm"
)

// User struct defines the user entity and their basic profile information. 👈 Struct Description

type User struct {
	// GORM Model Fields (Explicitly documented for Swagger)
	ID        uint           `json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2025-10-25T11:30:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2025-10-25T11:30:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Ignored in JSON output

	// User fields
	Username string `json:"username" gorm:"unique;not null" example:"user_alice"`     // Must be unique
	Email    string `json:"email" gorm:"unique;not null" example:"alice@example.com"` // Must be unique

	// Relationship: List of associated Todo items
	Todos []Todo `json:"todos"` // The 'json:"todos"' tag allows the list of todos to be included in the response.
}
