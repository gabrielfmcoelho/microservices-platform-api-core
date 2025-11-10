package domain

import (
	"context"

	"gorm.io/gorm"
)

// ContactIntent represents a contact request from the website
type ContactIntent struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Email       string `gorm:"size:255;not null"`
	Phone       string `gorm:"size:20;not null"`
	Company     string `gorm:"size:255;not null"`
	Message     string `gorm:"type:text"`
	ServiceName string `gorm:"size:255"`
	Status      string `gorm:"size:50;default:'pending'"` // pending, contacted, completed, cancelled
}

// CreateContactIntent represents the input for creating a new contact intent
type CreateContactIntent struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone" binding:"required"`
	Company     string `json:"company" binding:"required"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
}

// UpdateContactIntentStatus represents the input for updating contact intent status
type UpdateContactIntentStatus struct {
	Status string `json:"status" binding:"required,oneof=pending contacted completed cancelled"`
}

// PublicContactIntent represents the public view of a contact intent
type PublicContactIntent struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Company     string `json:"company"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// ContactIntentRepository defines the interface for contact intent data operations
type ContactIntentRepository interface {
	Create(ctx context.Context, contactIntent *ContactIntent) error
	Fetch(ctx context.Context) ([]ContactIntent, error)
	GetByID(ctx context.Context, id uint) (ContactIntent, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}

// ContactIntentUsecase defines the interface for contact intent business logic
type ContactIntentUsecase interface {
	Create(ctx context.Context, createContactIntent *CreateContactIntent) error
	Fetch(ctx context.Context) ([]PublicContactIntent, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
}
