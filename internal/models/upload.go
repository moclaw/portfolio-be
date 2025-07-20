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
	CreatedAt    time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// UploadResponse represents the response payload for upload operations
type UploadResponse struct {
	ID           uint      `json:"id" example:"1"`
	FileName     string    `json:"file_name" example:"image_123456.jpg"`
	OriginalName string    `json:"original_name" example:"my-image.jpg"`
	FileSize     int64     `json:"file_size" example:"1024000"`
	ContentType  string    `json:"content_type" example:"image/jpeg"`
	URL          string    `json:"url" example:"https://my-portfolio-bucket.s3.amazonaws.com/uploads/2023/01/01/image_123456.jpg"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

func (u *Upload) ToResponse() UploadResponse {
	return UploadResponse{
		ID:           u.ID,
		FileName:     u.FileName,
		OriginalName: u.OriginalName,
		FileSize:     u.FileSize,
		ContentType:  u.ContentType,
		URL:          u.URL,
		CreatedAt:    u.CreatedAt,
	}
}
