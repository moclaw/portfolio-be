package models

import (
	"time"

	"gorm.io/gorm"
)

// Content represents a content item in the system
type Content struct {
	ID          uint           `json:"id" gorm:"primarykey" example:"1"`
	Title       string         `json:"title" gorm:"not null" example:"My Blog Post"`
	Description string         `json:"description" example:"This is a sample blog post description"`
	Body        string         `json:"body" gorm:"type:text" example:"This is the content body of the blog post"`
	Category    string         `json:"category" example:"technology"`
	Tags        string         `json:"tags" example:"golang,api,backend"`
	Status      string         `json:"status" gorm:"default:draft" example:"published"`
	ImageURL    string         `json:"image_url" example:"https://example.com/image.jpg"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ContentRequest represents the request payload for creating/updating content
type ContentRequest struct {
	Title       string `json:"title" binding:"required" example:"My Blog Post"`
	Description string `json:"description" example:"This is a sample blog post description"`
	Body        string `json:"body" example:"This is the content body of the blog post"`
	Category    string `json:"category" example:"technology"`
	Tags        string `json:"tags" example:"golang,api,backend"`
	Status      string `json:"status" example:"published"`
	ImageURL    string `json:"image_url" example:"https://example.com/image.jpg"`
}

// ContentResponse represents the response payload for content operations
type ContentResponse struct {
	ID          uint      `json:"id" example:"1"`
	Title       string    `json:"title" example:"My Blog Post"`
	Description string    `json:"description" example:"This is a sample blog post description"`
	Body        string    `json:"body" example:"This is the content body of the blog post"`
	Category    string    `json:"category" example:"technology"`
	Tags        string    `json:"tags" example:"golang,api,backend"`
	Status      string    `json:"status" example:"published"`
	ImageURL    string    `json:"image_url" example:"https://example.com/image.jpg"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (c *Content) ToResponse() ContentResponse {
	return ContentResponse{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Body:        c.Body,
		Category:    c.Category,
		Tags:        c.Tags,
		Status:      c.Status,
		ImageURL:    c.ImageURL,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
