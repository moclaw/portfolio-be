package services

import (
	"fmt"
	"mime/multipart"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
	"strings"
	"time"
)

type UploadService struct {
	repo      *repository.UploadRepository
	s3Service *S3Service
}

func NewUploadService(repo *repository.UploadRepository, s3Service *S3Service) *UploadService {
	return &UploadService{
		repo:      repo,
		s3Service: s3Service,
	}
}

func (s *UploadService) UploadFile(file multipart.File, header *multipart.FileHeader) (*models.UploadResponse, error) {
	// Validate file type (optional - you can add more restrictions)
	allowedTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/svg+xml",
		"video/mp4", "video/webm", "video/ogg", "video/avi", "video/quicktime",
		"application/pdf", "text/plain", "application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Try to detect from filename
		ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])
		switch ext {
		case "jpg", "jpeg":
			contentType = "image/jpeg"
		case "png":
			contentType = "image/png"
		case "gif":
			contentType = "image/gif"
		case "webp":
			contentType = "image/webp"
		case "svg":
			contentType = "image/svg+xml"
		case "mp4":
			contentType = "video/mp4"
		case "webm":
			contentType = "video/webm"
		case "ogg":
			contentType = "video/ogg"
		case "avi":
			contentType = "video/avi"
		case "mov":
			contentType = "video/quicktime"
		case "pdf":
			contentType = "application/pdf"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Check if content type is allowed
	allowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("file type %s is not allowed", contentType)
	}

	// Upload to S3
	s3Key, url, err := s.s3Service.UploadFile(file, header)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Set expiry time for the URL (7 days from now)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Save to database
	upload := &models.Upload{
		FileName:     s3Key[strings.LastIndex(s3Key, "/")+1:], // Extract filename from s3 key
		OriginalName: header.Filename,
		FileSize:     header.Size,
		ContentType:  contentType,
		S3Key:        s3Key,
		S3Bucket:     s.s3Service.bucket,
		URL:          url,
		ExpiresAt:    &expiresAt,
		IsActive:     true,
	}

	if err := s.repo.Create(upload); err != nil {
		// If database save fails, try to clean up S3
		s.s3Service.DeleteFile(s3Key)
		return nil, fmt.Errorf("failed to save upload record: %w", err)
	}

	response := upload.ToResponse()
	return &response, nil
}

func (s *UploadService) GetUploadByID(id uint) (*models.UploadResponse, error) {
	upload, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload: %w", err)
	}

	response := upload.ToResponse()
	return &response, nil
}

func (s *UploadService) GetAllUploads(limit, offset int) ([]models.UploadResponse, error) {
	uploads, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get uploads: %w", err)
	}

	responses := make([]models.UploadResponse, len(uploads))
	for i, upload := range uploads {
		responses[i] = upload.ToResponse()
	}

	return responses, nil
}

func (s *UploadService) DeleteUpload(id uint) error {
	// Get upload record first
	upload, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get upload: %w", err)
	}

	// Delete from S3
	if err := s.s3Service.DeleteFile(upload.S3Key); err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	// Delete from database
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete upload record: %w", err)
	}

	return nil
}
