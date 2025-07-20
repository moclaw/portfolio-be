package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"1"`
	Name      string         `json:"name" gorm:"not null" example:"John Doe"`
	Email     string         `json:"email" gorm:"not null" example:"john@example.com"`
	Subject   string         `json:"subject" example:"Project Inquiry"`
	Message   string         `json:"message" gorm:"type:text;not null" example:"I would like to discuss a potential project."`
	Status    string         `json:"status" gorm:"default:unread" example:"unread"` // unread, read, replied
	IsActive  bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ContactRequest struct {
	Name    string `json:"name" binding:"required" example:"John Doe"`
	Email   string `json:"email" binding:"required,email" example:"john@example.com"`
	Subject string `json:"subject" example:"Project Inquiry"`
	Message string `json:"message" binding:"required" example:"I would like to discuss a potential project."`
}

type ContactResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Subject   string    `json:"subject" example:"Project Inquiry"`
	Message   string    `json:"message" example:"I would like to discuss a potential project."`
	Status    string    `json:"status" example:"unread"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type ContactUpdateRequest struct {
	Name     *string `json:"name,omitempty" example:"John Doe"`
	Email    *string `json:"email,omitempty" example:"john@example.com"`
	Subject  *string `json:"subject,omitempty" example:"Project Inquiry"`
	Message  *string `json:"message,omitempty" example:"I would like to discuss a potential project."`
	Status   *string `json:"status,omitempty" example:"read"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}
