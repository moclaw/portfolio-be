package models

import (
	"time"

	"gorm.io/gorm"
)

// Upload represents a file upload record in the system
type Upload struct {
	ID           uint           `json:"id" gorm:"primarykey" example:"1"`
	FileName     string         `json:"file_name" gorm:"not null" example:"image_123456.jpg"`
	OriginalName string         `json:"original_name" gorm:"not null" example:"my-image.jpg"`
	FileSize     int64          `json:"file_size" example:"1024000"`
	ContentType  string         `json:"content_type" example:"image/jpeg"`
	S3Key        string         `json:"s3_key" gorm:"not null;unique" example:"uploads/2023/01/01/image_123456.jpg"`
	S3Bucket     string         `json:"s3_bucket" gorm:"not null" example:"my-portfolio-bucket"`
	URL          string         `json:"url" gorm:"not null" example:"https://my-portfolio-bucket.s3.amazonaws.com/uploads/2023/01/01/image_123456.jpg"`
	ExpiresAt    *time.Time     `json:"expires_at" gorm:"index" example:"2024-01-01T00:00:00Z"`
	IsActive     bool           `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt    time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// UploadResponse represents the response payload for upload operations
type UploadResponse struct {
	ID           uint       `json:"id" example:"1"`
	FileName     string     `json:"file_name" example:"image_123456.jpg"`
	OriginalName string     `json:"original_name" example:"my-image.jpg"`
	FileSize     int64      `json:"file_size" example:"1024000"`
	ContentType  string     `json:"content_type" example:"image/jpeg"`
	URL          string     `json:"url" example:"https://my-portfolio-bucket.s3.amazonaws.com/uploads/2023/01/01/image_123456.jpg"`
	ExpiresAt    *time.Time `json:"expires_at" example:"2024-01-01T00:00:00Z"`
	IsActive     bool       `json:"is_active" example:"true"`
	CreatedAt    time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

func (u *Upload) ToResponse() UploadResponse {
	return UploadResponse{
		ID:           u.ID,
		FileName:     u.FileName,
		OriginalName: u.OriginalName,
		FileSize:     u.FileSize,
		ContentType:  u.ContentType,
		URL:          u.URL,
		ExpiresAt:    u.ExpiresAt,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
	}
}

// UploadSummary represents upload statistics
type UploadSummary struct {
	TotalFiles         int64  `json:"total_files" example:"150"`
	TotalSize          int64  `json:"total_size" example:"52428800"`
	TotalSizeFormatted string `json:"total_size_formatted" example:"50MB"`
	Images             int64  `json:"images" example:"80"`
	Documents          int64  `json:"documents" example:"30"`
	Videos             int64  `json:"videos" example:"20"`
	Others             int64  `json:"others" example:"20"`
}

// UploadListResponse represents the response for upload list with summary
type UploadListResponse struct {
	Uploads []UploadResponse `json:"uploads"`
	Summary UploadSummary    `json:"summary"`
}
