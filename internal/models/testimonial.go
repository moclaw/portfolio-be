package models

import (
	"time"

	"gorm.io/gorm"
)

// Testimonial represents a testimonial/recommendation in the portfolio
type Testimonial struct {
	ID          uint           `json:"id" gorm:"primarykey" example:"1"`
	Testimonial string         `json:"testimonial" gorm:"type:text;not null" example:"Great work! Highly recommended."`
	Name        string         `json:"name" gorm:"not null" example:"Sara Lee"`
	Designation string         `json:"designation" example:"CFO"`
	Company     string         `json:"company" example:"Acme Co"`
	Image       string         `json:"image" example:"https://randomuser.me/api/portraits/women/4.jpg"`
	Order       int            `json:"order" gorm:"column:sort_order;default:0" example:"1"`
	IsActive    bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TestimonialRequest represents the request payload for creating/updating testimonial
type TestimonialRequest struct {
	Testimonial string `json:"testimonial" binding:"required" example:"Great work! Highly recommended."`
	Name        string `json:"name" binding:"required" example:"Sara Lee"`
	Designation string `json:"designation" example:"CFO"`
	Company     string `json:"company" example:"Acme Co"`
	Image       string `json:"image" example:"https://randomuser.me/api/portraits/women/4.jpg"`
	Order       int    `json:"order" example:"1"`
	IsActive    bool   `json:"is_active" example:"true"`
}

// TestimonialResponse represents the response payload for testimonial operations
type TestimonialResponse struct {
	ID          uint      `json:"id" example:"1"`
	Testimonial string    `json:"testimonial" example:"Great work! Highly recommended."`
	Name        string    `json:"name" example:"Sara Lee"`
	Designation string    `json:"designation" example:"CFO"`
	Company     string    `json:"company" example:"Acme Co"`
	Image       string    `json:"image" example:"https://randomuser.me/api/portraits/women/4.jpg"`
	Order       int       `json:"order" example:"1"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (t *Testimonial) ToResponse() TestimonialResponse {
	return TestimonialResponse{
		ID:          t.ID,
		Testimonial: t.Testimonial,
		Name:        t.Name,
		Designation: t.Designation,
		Company:     t.Company,
		Image:       t.Image,
		Order:       t.Order,
		IsActive:    t.IsActive,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
