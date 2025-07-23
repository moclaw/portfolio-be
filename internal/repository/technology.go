package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type TechnologyRepository interface {
	Create(technology *models.Technology) error
	GetAll() ([]models.Technology, error)
	GetByID(id uint) (*models.Technology, error)
	Update(technology *models.Technology) error
	Delete(id uint) error
	GetActive() ([]models.Technology, error)
	GetByCategory(category string) ([]models.Technology, error)
	GetCount() (int64, error)
	UpdateOrder(id uint, order int) error
}

type technologyRepository struct {
	db *gorm.DB
}

func NewTechnologyRepository(db *gorm.DB) TechnologyRepository {
	return &technologyRepository{db: db}
}

func (r *technologyRepository) Create(technology *models.Technology) error {
	return r.db.Create(technology).Error
}

func (r *technologyRepository) GetAll() ([]models.Technology, error) {
	var technologies []models.Technology
	err := r.db.Find(&technologies).Error
	return technologies, err
}

func (r *technologyRepository) GetByID(id uint) (*models.Technology, error) {
	var technology models.Technology
	err := r.db.First(&technology, id).Error
	if err != nil {
		return nil, err
	}
	return &technology, nil
}

func (r *technologyRepository) Update(technology *models.Technology) error {
	return r.db.Save(technology).Error
}

func (r *technologyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Technology{}, id).Error
}

func (r *technologyRepository) GetActive() ([]models.Technology, error) {
	var technologies []models.Technology
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, created_at ASC").Find(&technologies).Error
	return technologies, err
}

func (r *technologyRepository) GetByCategory(category string) ([]models.Technology, error) {
	var technologies []models.Technology
	err := r.db.Where("category = ? AND is_active = ?", category, true).Order("sort_order ASC, created_at ASC").Find(&technologies).Error
	return technologies, err
}

func (r *technologyRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Technology{}).Count(&count).Error
	return count, err
}

func (r *technologyRepository) UpdateOrder(id uint, order int) error {
	return r.db.Model(&models.Technology{}).Where("id = ?", id).Update("sort_order", order).Error
}
