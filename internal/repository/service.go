package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *models.Service) error
	GetAll() ([]models.Service, error)
	GetByID(id uint) (*models.Service, error)
	Update(service *models.Service) error
	Delete(id uint) error
	GetActive() ([]models.Service, error)
	GetCount() (int64, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *models.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) GetAll() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&services).Error
	return services, err
}

func (r *serviceRepository) GetByID(id uint) (*models.Service, error) {
	var service models.Service
	err := r.db.First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) Update(service *models.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Service{}, id).Error
}

func (r *serviceRepository) GetActive() ([]models.Service, error) {
	var services []models.Service
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, created_at ASC").Find(&services).Error
	return services, err
}

func (r *serviceRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Service{}).Count(&count).Error
	return count, err
}
