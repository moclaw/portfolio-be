package services

import (
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type ContactService struct {
	contactRepo *repository.ContactRepository
}

func NewContactService(contactRepo *repository.ContactRepository) *ContactService {
	return &ContactService{
		contactRepo: contactRepo,
	}
}

func (s *ContactService) CreateContact(req *models.ContactRequest) (*models.Contact, error) {
	contact := &models.Contact{
		Name:     req.Name,
		Email:    req.Email,
		Subject:  req.Subject,
		Message:  req.Message,
		Status:   "unread",
		IsActive: true,
	}

	if err := s.contactRepo.Create(contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactService) GetAllContacts(page, limit int) ([]models.ContactResponse, int64, error) {
	contacts, total, err := s.contactRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.ContactResponse
	for _, contact := range contacts {
		responses = append(responses, models.ContactResponse{
			ID:        contact.ID,
			Name:      contact.Name,
			Email:     contact.Email,
			Subject:   contact.Subject,
			Message:   contact.Message,
			Status:    contact.Status,
			IsActive:  contact.IsActive,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		})
	}

	return responses, total, nil
}

func (s *ContactService) GetContactByID(id uint) (*models.ContactResponse, error) {
	contact, err := s.contactRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &models.ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Subject:   contact.Subject,
		Message:   contact.Message,
		Status:    contact.Status,
		IsActive:  contact.IsActive,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	return response, nil
}

func (s *ContactService) UpdateContact(id uint, req *models.ContactUpdateRequest) (*models.ContactResponse, error) {
	contact, err := s.contactRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Email != nil {
		contact.Email = *req.Email
	}
	if req.Subject != nil {
		contact.Subject = *req.Subject
	}
	if req.Message != nil {
		contact.Message = *req.Message
	}
	if req.Status != nil {
		contact.Status = *req.Status
	}
	if req.IsActive != nil {
		contact.IsActive = *req.IsActive
	}

	if err := s.contactRepo.Update(contact); err != nil {
		return nil, err
	}

	response := &models.ContactResponse{
		ID:        contact.ID,
		Name:      contact.Name,
		Email:     contact.Email,
		Subject:   contact.Subject,
		Message:   contact.Message,
		Status:    contact.Status,
		IsActive:  contact.IsActive,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	return response, nil
}

func (s *ContactService) DeleteContact(id uint) error {
	return s.contactRepo.Delete(id)
}

func (s *ContactService) GetContactsByStatus(status string, page, limit int) ([]models.ContactResponse, int64, error) {
	contacts, total, err := s.contactRepo.GetByStatus(status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []models.ContactResponse
	for _, contact := range contacts {
		responses = append(responses, models.ContactResponse{
			ID:        contact.ID,
			Name:      contact.Name,
			Email:     contact.Email,
			Subject:   contact.Subject,
			Message:   contact.Message,
			Status:    contact.Status,
			IsActive:  contact.IsActive,
			CreatedAt: contact.CreatedAt,
			UpdatedAt: contact.UpdatedAt,
		})
	}

	return responses, total, nil
}

func (s *ContactService) MarkAsRead(id uint) error {
	return s.contactRepo.MarkAsRead(id)
}

func (s *ContactService) GetUnreadCount() (int64, error) {
	return s.contactRepo.GetUnreadCount()
}

func (s *ContactService) GetContactsCount() (int64, error) {
	return s.contactRepo.GetCount()
}
