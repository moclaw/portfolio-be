package models

import (
	"time"

	"gorm.io/gorm"
)

// Service represents a service offered in the portfolio
type Service struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"1"`
	Title     string         `json:"title" gorm:"not null" example:"Full Stack Developer"`
	Icon      string         `json:"icon" example:"web.png"`
	Order     int            `json:"order" gorm:"column:sort_order;default:0" example:"1"`
	IsActive  bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ServiceRequest represents the request payload for creating/updating service
type ServiceRequest struct {
	Title    string `json:"title" binding:"required" example:"Full Stack Developer"`
	Icon     string `json:"icon" example:"web.png"`
	Order    int    `json:"order" example:"1"`
	IsActive bool   `json:"is_active" example:"true"`
}

// ServiceResponse represents the response payload for service operations
type ServiceResponse struct {
	ID        uint      `json:"id" example:"1"`
	Title     string    `json:"title" example:"Full Stack Developer"`
	Icon      string    `json:"icon" example:"web.png"`
	Order     int       `json:"order" example:"1"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (s *Service) ToResponse() ServiceResponse {
	return ServiceResponse{
		ID:        s.ID,
		Title:     s.Title,
		Icon:      s.Icon,
		Order:     s.Order,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
