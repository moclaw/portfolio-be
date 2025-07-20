package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Project represents a project in the portfolio
type Project struct {
	ID             uint           `json:"id" gorm:"primarykey" example:"1"`
	Name           string         `json:"name" gorm:"not null" example:"Car Rent"`
	Description    string         `json:"description" gorm:"type:text" example:"Web-based platform for car rentals"`
	Tags           string         `json:"tags" gorm:"type:text" example:"[{\"name\":\"react\",\"color\":\"blue-text-gradient\"}]"`
	Image          string         `json:"image" example:"carrent.png"`
	SourceCodeLink string         `json:"source_code_link" example:"https://github.com/example"`
	LiveDemoLink   string         `json:"live_demo_link" example:"https://example.com"`
	Order          int            `json:"order" gorm:"column:sort_order;default:0" example:"1"`
	IsActive       bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt      time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// ProjectTag represents a tag for a project
type ProjectTag struct {
	Name  string `json:"name" example:"react"`
	Color string `json:"color" example:"blue-text-gradient"`
}

// ProjectRequest represents the request payload for creating/updating project
type ProjectRequest struct {
	Name           string       `json:"name" binding:"required" example:"Car Rent"`
	Description    string       `json:"description" example:"Web-based platform for car rentals"`
	Tags           []ProjectTag `json:"tags"`
	Image          string       `json:"image" example:"carrent.png"`
	SourceCodeLink string       `json:"source_code_link" example:"https://github.com/example"`
	LiveDemoLink   string       `json:"live_demo_link" example:"https://example.com"`
	Order          int          `json:"order" example:"1"`
	IsActive       bool         `json:"is_active" example:"true"`
}

// ProjectResponse represents the response payload for project operations
type ProjectResponse struct {
	ID             uint         `json:"id" example:"1"`
	Name           string       `json:"name" example:"Car Rent"`
	Description    string       `json:"description" example:"Web-based platform for car rentals"`
	Tags           []ProjectTag `json:"tags"`
	Image          string       `json:"image" example:"carrent.png"`
	SourceCodeLink string       `json:"source_code_link" example:"https://github.com/example"`
	LiveDemoLink   string       `json:"live_demo_link" example:"https://example.com"`
	Order          int          `json:"order" example:"1"`
	IsActive       bool         `json:"is_active" example:"true"`
	CreatedAt      time.Time    `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time    `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (p *Project) ToResponse() ProjectResponse {
	// Parse JSON string back to slice
	var tags []ProjectTag
	if p.Tags != "" {
		json.Unmarshal([]byte(p.Tags), &tags)
	}

	return ProjectResponse{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Tags:           tags,
		Image:          p.Image,
		SourceCodeLink: p.SourceCodeLink,
		LiveDemoLink:   p.LiveDemoLink,
		Order:          p.Order,
		IsActive:       p.IsActive,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
