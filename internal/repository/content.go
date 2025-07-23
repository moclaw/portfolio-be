package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) Create(content *models.Content) error {
	return r.db.Create(content).Error
}

func (r *ContentRepository) GetByID(id uint) (*models.Content, error) {
	var content models.Content
	err := r.db.First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *ContentRepository) GetAll(limit, offset int) ([]models.Content, error) {
	var contents []models.Content
	err := r.db.Limit(limit).Offset(offset).Find(&contents).Error
	return contents, err
}

func (r *ContentRepository) GetByCategory(category string, limit, offset int) ([]models.Content, error) {
	var contents []models.Content
	err := r.db.Where("category = ?", category).Limit(limit).Offset(offset).Find(&contents).Error
	return contents, err
}

func (r *ContentRepository) GetByStatus(status string, limit, offset int) ([]models.Content, error) {
	var contents []models.Content
	err := r.db.Where("status = ?", status).Limit(limit).Offset(offset).Find(&contents).Error
	return contents, err
}

func (r *ContentRepository) Update(content *models.Content) error {
	return r.db.Save(content).Error
}

func (r *ContentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Content{}, id).Error
}

func (r *ContentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Content{}).Count(&count).Error
	return count, err
}
