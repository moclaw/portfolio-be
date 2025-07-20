package repository

import (
	"portfolio-be/internal/models"

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

func (r *UploadRepository) Delete(id uint) error {
	return r.db.Delete(&models.Upload{}, id).Error
}

func (r *UploadRepository) DeleteByS3Key(s3Key string) error {
	return r.db.Where("s3_key = ?", s3Key).Delete(&models.Upload{}).Error
}