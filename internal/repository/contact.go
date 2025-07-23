package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(contact *models.Contact) error {
	return r.db.Create(contact).Error
}

func (r *ContactRepository) GetAll(page, limit int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	// Count total records
	if err := r.db.Model(&models.Contact{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&contacts).Error

	return contacts, total, err
}

func (r *ContactRepository) GetByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, id).Error
	return &contact, err
}

func (r *ContactRepository) Update(contact *models.Contact) error {
	return r.db.Save(contact).Error
}

func (r *ContactRepository) Delete(id uint) error {
	return r.db.Delete(&models.Contact{}, id).Error
}

func (r *ContactRepository) GetByStatus(status string, page, limit int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	query := r.db.Model(&models.Contact{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&contacts).Error

	return contacts, total, err
}

func (r *ContactRepository) MarkAsRead(id uint) error {
	return r.db.Model(&models.Contact{}).
		Where("id = ?", id).
		Update("status", "read").Error
}

func (r *ContactRepository) GetUnreadCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Contact{}).
		Where("status = ?", "unread").
		Count(&count).Error
	return count, err
}

func (r *ContactRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Contact{}).Count(&count).Error
	return count, err
}
