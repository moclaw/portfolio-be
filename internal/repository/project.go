package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	GetAll() ([]models.Project, error)
	GetByID(id uint) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id uint) error
	GetActive() ([]models.Project, error)
	GetCount() (int64, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

func (r *projectRepository) GetActive() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, created_at ASC").Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Project{}).Count(&count).Error
	return count, err
}
