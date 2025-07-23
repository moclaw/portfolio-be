package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type ExperienceRepository interface {
	Create(experience *models.Experience) error
	GetAll() ([]models.Experience, error)
	GetByID(id uint) (*models.Experience, error)
	Update(experience *models.Experience) error
	Delete(id uint) error
	GetActive() ([]models.Experience, error)
	GetCount() (int64, error)
}

type experienceRepository struct {
	db *gorm.DB
}

func NewExperienceRepository(db *gorm.DB) ExperienceRepository {
	return &experienceRepository{db: db}
}

func (r *experienceRepository) Create(experience *models.Experience) error {
	return r.db.Create(experience).Error
}

func (r *experienceRepository) GetAll() ([]models.Experience, error) {
	var experiences []models.Experience
	err := r.db.Order("order_index ASC, created_at ASC").Find(&experiences).Error
	return experiences, err
}

func (r *experienceRepository) GetByID(id uint) (*models.Experience, error) {
	var experience models.Experience
	err := r.db.First(&experience, id).Error
	if err != nil {
		return nil, err
	}
	return &experience, nil
}

func (r *experienceRepository) Update(experience *models.Experience) error {
	return r.db.Save(experience).Error
}

func (r *experienceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Experience{}, id).Error
}

func (r *experienceRepository) GetActive() ([]models.Experience, error) {
	var experiences []models.Experience
	err := r.db.Where("is_active = ?", true).Order("order_index ASC, created_at ASC").Find(&experiences).Error
	return experiences, err
}

func (r *experienceRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Experience{}).Count(&count).Error
	return count, err
}
