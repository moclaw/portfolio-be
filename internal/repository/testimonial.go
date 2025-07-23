package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type TestimonialRepository interface {
	Create(testimonial *models.Testimonial) error
	GetAll() ([]models.Testimonial, error)
	GetByID(id uint) (*models.Testimonial, error)
	Update(testimonial *models.Testimonial) error
	Delete(id uint) error
	GetActive() ([]models.Testimonial, error)
	GetCount() (int64, error)
}

type testimonialRepository struct {
	db *gorm.DB
}

func NewTestimonialRepository(db *gorm.DB) TestimonialRepository {
	return &testimonialRepository{db: db}
}

func (r *testimonialRepository) Create(testimonial *models.Testimonial) error {
	return r.db.Create(testimonial).Error
}

func (r *testimonialRepository) GetAll() ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&testimonials).Error
	return testimonials, err
}

func (r *testimonialRepository) GetByID(id uint) (*models.Testimonial, error) {
	var testimonial models.Testimonial
	err := r.db.First(&testimonial, id).Error
	if err != nil {
		return nil, err
	}
	return &testimonial, nil
}

func (r *testimonialRepository) Update(testimonial *models.Testimonial) error {
	return r.db.Save(testimonial).Error
}

func (r *testimonialRepository) Delete(id uint) error {
	return r.db.Delete(&models.Testimonial{}, id).Error
}

func (r *testimonialRepository) GetActive() ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	err := r.db.Where("is_active = ?", true).Order("sort_order ASC, created_at ASC").Find(&testimonials).Error
	return testimonials, err
}

func (r *testimonialRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Testimonial{}).Count(&count).Error
	return count, err
}
