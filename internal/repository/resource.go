package repository

import (
	"portfolio-be/internal/models"
	"time"

	"gorm.io/gorm"
)

type ResourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

func (r *ResourceRepository) Create(resource *models.Resource) error {
	return r.db.Create(resource).Error
}

func (r *ResourceRepository) GetByID(id uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Preload("Upload").First(&resource, id).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *ResourceRepository) GetAll(limit, offset int) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Preload("Upload").Limit(limit).Offset(offset).Find(&resources).Error
	return resources, err
}

func (r *ResourceRepository) GetByType(resourceType models.ResourceType, limit, offset int) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Preload("Upload").Where("type = ?", resourceType).Limit(limit).Offset(offset).Find(&resources).Error
	return resources, err
}

func (r *ResourceRepository) GetByCategory(category string, limit, offset int) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Preload("Upload").Where("category = ?", category).Limit(limit).Offset(offset).Find(&resources).Error
	return resources, err
}

func (r *ResourceRepository) GetPublic(limit, offset int) ([]models.Resource, error) {
	var resources []models.Resource
	err := r.db.Preload("Upload").Where("is_public = ? AND is_active = ?", true, true).Limit(limit).Offset(offset).Find(&resources).Error
	return resources, err
}

func (r *ResourceRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Resource{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ResourceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Resource{}, id).Error
}

func (r *ResourceRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Resource{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *ResourceRepository) IncrementDownloadCount(id uint) error {
	return r.db.Model(&models.Resource{}).Where("id = ?", id).Update("download_count", gorm.Expr("download_count + 1")).Error
}

func (r *ResourceRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Resource{}).Count(&count).Error
	return count, err
}

func (r *ResourceRepository) CountByType(resourceType models.ResourceType) (int64, error) {
	var count int64
	err := r.db.Model(&models.Resource{}).Where("type = ?", resourceType).Count(&count).Error
	return count, err
}

func (r *ResourceRepository) CountByCategory(category string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Resource{}).Where("category = ?", category).Count(&count).Error
	return count, err
}

// GetExpiringSoon returns resources whose uploads are expiring within the specified duration
func (r *ResourceRepository) GetExpiringSoon(duration time.Duration) ([]models.Resource, error) {
	var resources []models.Resource
	expiryThreshold := time.Now().Add(duration)

	err := r.db.Preload("Upload").
		Joins("JOIN uploads ON resources.upload_id = uploads.id").
		Where("uploads.expires_at IS NOT NULL AND uploads.expires_at <= ? AND uploads.expires_at > ?", expiryThreshold, time.Now()).
		Find(&resources).Error

	return resources, err
}

// GetExpired returns resources whose uploads have expired
func (r *ResourceRepository) GetExpired() ([]models.Resource, error) {
	var resources []models.Resource
	now := time.Now()

	err := r.db.Preload("Upload").
		Joins("JOIN uploads ON resources.upload_id = uploads.id").
		Where("uploads.expires_at IS NOT NULL AND uploads.expires_at <= ?", now).
		Find(&resources).Error

	return resources, err
}

// Search resources by name, description, or tags
func (r *ResourceRepository) Search(query string, limit, offset int) ([]models.Resource, error) {
	var resources []models.Resource
	searchPattern := "%" + query + "%"

	err := r.db.Preload("Upload").
		Where("name LIKE ? OR description LIKE ? OR tags LIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).Offset(offset).
		Find(&resources).Error

	return resources, err
}
