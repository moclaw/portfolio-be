package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Experience represents work experience in the portfolio
type Experience struct {
	ID          uint           `json:"id" gorm:"primarykey" example:"1"`
	Title       string         `json:"title" gorm:"not null" example:"Full Stack Developer"`
	CompanyName string         `json:"company_name" gorm:"not null" example:"MaicoGroup"`
	Icon        string         `json:"icon" example:"maico.png"`
	IconBg      string         `json:"icon_bg" example:"#ffffff"`
	Date        string         `json:"date" gorm:"not null" example:"Apr 2021 - Jan 2022"`
	Points      string         `json:"points" gorm:"type:text" example:"[\"Point 1\", \"Point 2\"]"`
	Order       int            `json:"order" gorm:"column:order_index;default:0" example:"1"`
	IsActive    bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ExperienceRequest represents the request payload for creating/updating experience
type ExperienceRequest struct {
	Title       string   `json:"title" binding:"required" example:"Full Stack Developer"`
	CompanyName string   `json:"company_name" binding:"required" example:"MaicoGroup"`
	Icon        string   `json:"icon" example:"maico.png"`
	IconBg      string   `json:"icon_bg" example:"#ffffff"`
	Date        string   `json:"date" binding:"required" example:"Apr 2021 - Jan 2022"`
	Points      []string `json:"points" example:"[\"Point 1\", \"Point 2\"]"`
	Order       int      `json:"order" example:"1"`
	IsActive    bool     `json:"is_active" example:"true"`
}

// ExperienceResponse represents the response payload for experience operations
type ExperienceResponse struct {
	ID          uint      `json:"id" example:"1"`
	Title       string    `json:"title" example:"Full Stack Developer"`
	CompanyName string    `json:"company_name" example:"MaicoGroup"`
	Icon        string    `json:"icon" example:"maico.png"`
	IconBg      string    `json:"icon_bg" example:"#ffffff"`
	Date        string    `json:"date" example:"Apr 2021 - Jan 2022"`
	Points      []string  `json:"points" example:"[\"Point 1\", \"Point 2\"]"`
	Order       int       `json:"order" example:"1"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func (e *Experience) ToResponse() ExperienceResponse {
	// Parse JSON string back to slice
	var points []string
	if e.Points != "" {
		json.Unmarshal([]byte(e.Points), &points)
	}

	return ExperienceResponse{
		ID:          e.ID,
		Title:       e.Title,
		CompanyName: e.CompanyName,
		Icon:        e.Icon,
		IconBg:      e.IconBg,
		Date:        e.Date,
		Points:      points,
		Order:       e.Order,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
