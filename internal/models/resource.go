package models

import (
	"time"

	"gorm.io/gorm"
)

// ResourceType represents the type of resource
type ResourceType string

const (
	ResourceTypeImage    ResourceType = "image"
	ResourceTypeDocument ResourceType = "document"
	ResourceTypeVideo    ResourceType = "video"
	ResourceTypeAudio    ResourceType = "audio"
	ResourceTypeOther    ResourceType = "other"
)

// Resource represents a digital resource/asset in the system
type Resource struct {
	ID            uint           `json:"id" gorm:"primarykey" example:"1"`
	Name          string         `json:"name" gorm:"not null" example:"Portfolio Banner"`
	Description   string         `json:"description" example:"Main banner image for portfolio"`
	Type          ResourceType   `json:"type" gorm:"not null" example:"image"`
	Category      string         `json:"category" example:"banner"`
	Tags          string         `json:"tags" example:"portfolio,banner,hero"`
	UploadID      uint           `json:"upload_id" gorm:"not null" example:"1"`
	Upload        Upload         `json:"upload" gorm:"foreignKey:UploadID"`
	Alt           string         `json:"alt" example:"Portfolio banner showing professional work"`
	IsPublic      bool           `json:"is_public" gorm:"default:true" example:"true"`
	IsActive      bool           `json:"is_active" gorm:"default:true" example:"true"`
	ViewCount     int64          `json:"view_count" gorm:"default:0" example:"0"`
	DownloadCount int64          `json:"download_count" gorm:"default:0" example:"0"`
	CreatedAt     time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// ResourceResponse represents the response payload for resource operations
type ResourceResponse struct {
	ID            uint           `json:"id" example:"1"`
	Name          string         `json:"name" example:"Portfolio Banner"`
	Description   string         `json:"description" example:"Main banner image for portfolio"`
	Type          ResourceType   `json:"type" example:"image"`
	Category      string         `json:"category" example:"banner"`
	Tags          []string       `json:"tags" example:"portfolio,banner,hero"`
	Upload        UploadResponse `json:"upload"`
	Alt           string         `json:"alt" example:"Portfolio banner showing professional work"`
	IsPublic      bool           `json:"is_public" example:"true"`
	IsActive      bool           `json:"is_active" example:"true"`
	ViewCount     int64          `json:"view_count" example:"0"`
	DownloadCount int64          `json:"download_count" example:"0"`
	CreatedAt     time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ResourceCreateRequest represents the request payload for creating a resource
type ResourceCreateRequest struct {
	Name        string       `json:"name" binding:"required" example:"Portfolio Banner"`
	Description string       `json:"description" example:"Main banner image for portfolio"`
	Type        ResourceType `json:"type" binding:"required" example:"image"`
	Category    string       `json:"category" example:"banner"`
	Tags        string       `json:"tags" example:"portfolio,banner,hero"`
	UploadID    uint         `json:"upload_id" binding:"required" example:"1"`
	Alt         string       `json:"alt" example:"Portfolio banner showing professional work"`
	IsPublic    *bool        `json:"is_public" example:"true"`
	IsActive    *bool        `json:"is_active" example:"true"`
}

// ResourceUpdateRequest represents the request payload for updating a resource
type ResourceUpdateRequest struct {
	Name        *string       `json:"name,omitempty" example:"Portfolio Banner Updated"`
	Description *string       `json:"description,omitempty" example:"Updated description"`
	Type        *ResourceType `json:"type,omitempty" example:"image"`
	Category    *string       `json:"category,omitempty" example:"banner"`
	Tags        *string       `json:"tags,omitempty" example:"portfolio,banner,hero,updated"`
	Alt         *string       `json:"alt,omitempty" example:"Updated alt text"`
	IsPublic    *bool         `json:"is_public,omitempty" example:"true"`
	IsActive    *bool         `json:"is_active,omitempty" example:"true"`
}

func (r *Resource) ToResponse() ResourceResponse {
	tags := []string{}
	if r.Tags != "" {
		// Simple split by comma - you might want to use a more sophisticated parser
		tags = append(tags, r.Tags)
	}

	return ResourceResponse{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description,
		Type:          r.Type,
		Category:      r.Category,
		Tags:          tags,
		Upload:        r.Upload.ToResponse(),
		Alt:           r.Alt,
		IsPublic:      r.IsPublic,
		IsActive:      r.IsActive,
		ViewCount:     r.ViewCount,
		DownloadCount: r.DownloadCount,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

// IsExpiringSoon checks if the associated upload is expiring within the specified duration
func (r *Resource) IsExpiringSoon(duration time.Duration) bool {
	if r.Upload.ExpiresAt == nil {
		return false
	}
	return time.Until(*r.Upload.ExpiresAt) <= duration
}

// IsExpired checks if the associated upload has expired
func (r *Resource) IsExpired() bool {
	if r.Upload.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.Upload.ExpiresAt)
}
