package repository

import (
	"fmt"
	"portfolio-be/internal/models"
	"time"

	"gorm.io/gorm"
)

type UploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) Create(upload *models.Upload) error {
	return r.db.Create(upload).Error
}

func (r *UploadRepository) GetByID(id uint) (*models.Upload, error) {
	var upload models.Upload
	err := r.db.First(&upload, id).Error
	if err != nil {
		return nil, err
	}
	return &upload, nil
}

func (r *UploadRepository) GetByS3Key(s3Key string) (*models.Upload, error) {
	var upload models.Upload
	err := r.db.Where("s3_key = ?", s3Key).First(&upload).Error
	if err != nil {
		return nil, err
	}
	return &upload, nil
}

func (r *UploadRepository) GetAll(limit, offset int) ([]models.Upload, error) {
	var uploads []models.Upload
	err := r.db.Limit(limit).Offset(offset).Find(&uploads).Error
	return uploads, err
}

func (r *UploadRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Upload{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *UploadRepository) Delete(id uint) error {
	return r.db.Delete(&models.Upload{}, id).Error
}

func (r *UploadRepository) DeleteByS3Key(s3Key string) error {
	return r.db.Where("s3_key = ?", s3Key).Delete(&models.Upload{}).Error
}

// UpdateURL updates the URL of an upload record
func (r *UploadRepository) UpdateURL(id uint, newURL string) error {
	return r.db.Model(&models.Upload{}).Where("id = ?", id).Update("url", newURL).Error
}

// UpdateExpiry updates the expiry time of an upload record
func (r *UploadRepository) UpdateExpiry(id uint, expiresAt *time.Time) error {
	return r.db.Model(&models.Upload{}).Where("id = ?", id).Update("expires_at", expiresAt).Error
}

// GetExpiringSoon returns uploads that are expiring within the specified duration
func (r *UploadRepository) GetExpiringSoon(duration time.Duration) ([]models.Upload, error) {
	var uploads []models.Upload
	expiryThreshold := time.Now().Add(duration)

	err := r.db.Where("expires_at IS NOT NULL AND expires_at <= ? AND expires_at > ?", expiryThreshold, time.Now()).Find(&uploads).Error
	return uploads, err
}

// GetExpired returns uploads that have expired
func (r *UploadRepository) GetExpired() ([]models.Upload, error) {
	var uploads []models.Upload
	now := time.Now()

	err := r.db.Where("expires_at IS NOT NULL AND expires_at <= ?", now).Find(&uploads).Error
	return uploads, err
}

// GetUploadSummary returns upload statistics
func (r *UploadRepository) GetUploadSummary() (*models.UploadSummary, error) {
	var summary models.UploadSummary

	// Get total files count
	err := r.db.Model(&models.Upload{}).Where("is_active = ?", true).Count(&summary.TotalFiles).Error
	if err != nil {
		return nil, err
	}

	// Get total size
	var totalSize struct {
		Total int64
	}
	err = r.db.Model(&models.Upload{}).Where("is_active = ?", true).Select("COALESCE(SUM(file_size), 0) as total").Scan(&totalSize).Error
	if err != nil {
		return nil, err
	}
	summary.TotalSize = totalSize.Total
	summary.TotalSizeFormatted = formatFileSize(totalSize.Total)

	// Get images count
	err = r.db.Model(&models.Upload{}).Where("is_active = ? AND content_type LIKE ?", true, "image/%").Count(&summary.Images).Error
	if err != nil {
		return nil, err
	}

	// Get documents count
	err = r.db.Model(&models.Upload{}).Where("is_active = ? AND (content_type LIKE ? OR content_type LIKE ? OR content_type = ?)",
		true, "application/pdf", "application/%word%", "text/plain").Count(&summary.Documents).Error
	if err != nil {
		return nil, err
	}

	// Get videos count
	err = r.db.Model(&models.Upload{}).Where("is_active = ? AND content_type LIKE ?", true, "video/%").Count(&summary.Videos).Error
	if err != nil {
		return nil, err
	}

	// Calculate others
	summary.Others = summary.TotalFiles - summary.Images - summary.Documents - summary.Videos

	return &summary, nil
}

// formatFileSize formats file size in bytes to human readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
