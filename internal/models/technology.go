package models

import (
	"time"

	"gorm.io/gorm"
)

// Technology represents a technology skill in the portfolio
type Technology struct {
	ID        uint           `json:"id" gorm:"primarykey" example:"1"`
	Name      string         `json:"name" gorm:"not null" example:"C#"`
	Icon      string         `json:"icon" example:"csharp.png"`
	Category  string         `json:"category" example:"programming"`
	Order     int            `json:"order" gorm:"column:sort_order;default:0" example:"1"`
	IsActive  bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TechnologyRequest represents the request payload for creating/updating technology
type TechnologyRequest struct {
	Name     string `json:"name" binding:"required" example:"C#"`
	Icon     string `json:"icon" example:"csharp.png"`
	Category string `json:"category" example:"programming"`
	Order    int    `json:"order" example:"1"`
	IsActive bool   `json:"is_active" example:"true"`
}

// TechnologyResponse represents the response payload for technology operations
type TechnologyResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"C#"`
	Icon      string    `json:"icon" example:"csharp.png"`
	Category  string    `json:"category" example:"programming"`
	Order     int       `json:"order" example:"1"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (t *Technology) ToResponse() TechnologyResponse {
	return TechnologyResponse{
		ID:        t.ID,
		Name:      t.Name,
		Icon:      t.Icon,
		Category:  t.Category,
		Order:     t.Order,
		IsActive:  t.IsActive,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
