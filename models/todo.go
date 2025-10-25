package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo struct defines a single task item in the system 👈 ADD THIS LINE
type Todo struct {
	// GORM fields explicitly documented for Swagger
	ID        uint           `json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2025-10-25T10:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2025-10-25T10:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Todo fields
	Item      string `json:"item" gorm:"not null" example:"Buy groceries"`
	Completed bool   `json:"completed" example:"false"`
	UserID    uint   `json:"user_id" example:"1"` // Foreign key linking to User.ID
}
