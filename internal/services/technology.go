package services

import (
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type TechnologyService interface {
	CreateTechnology(request *models.TechnologyRequest) (*models.TechnologyResponse, error)
	GetAllTechnologies() ([]models.TechnologyResponse, error)
	GetTechnologyByID(id uint) (*models.TechnologyResponse, error)
	UpdateTechnology(id uint, request *models.TechnologyRequest) (*models.TechnologyResponse, error)
	DeleteTechnology(id uint) error
	GetActiveTechnologies() ([]models.TechnologyResponse, error)
	GetTechnologiesByCategory(category string) ([]models.TechnologyResponse, error)
	GetTechnologiesCount() (int64, error)
	UpdateTechnologiesOrder(items []OrderItem) error
}

type technologyService struct {
	technologyRepo repository.TechnologyRepository
}

func NewTechnologyService(technologyRepo repository.TechnologyRepository) TechnologyService {
	return &technologyService{technologyRepo: technologyRepo}
}

func (s *technologyService) CreateTechnology(request *models.TechnologyRequest) (*models.TechnologyResponse, error) {
	technology := &models.Technology{
		Name:     request.Name,
		Icon:     request.Icon,
		Category: request.Category,
		Order:    request.Order,
		IsActive: request.IsActive,
	}

	err := s.technologyRepo.Create(technology)
	if err != nil {
		return nil, err
	}

	response := technology.ToResponse()
	return &response, nil
}

func (s *technologyService) GetAllTechnologies() ([]models.TechnologyResponse, error) {
	technologies, err := s.technologyRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.TechnologyResponse
	for _, technology := range technologies {
		responses = append(responses, technology.ToResponse())
	}

	return responses, nil
}

func (s *technologyService) GetTechnologyByID(id uint) (*models.TechnologyResponse, error) {
	technology, err := s.technologyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := technology.ToResponse()
	return &response, nil
}

func (s *technologyService) UpdateTechnology(id uint, request *models.TechnologyRequest) (*models.TechnologyResponse, error) {
	technology, err := s.technologyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	technology.Name = request.Name
	technology.Icon = request.Icon
	technology.Category = request.Category
	technology.Order = request.Order
	technology.IsActive = request.IsActive

	err = s.technologyRepo.Update(technology)
	if err != nil {
		return nil, err
	}

	response := technology.ToResponse()
	return &response, nil
}

func (s *technologyService) DeleteTechnology(id uint) error {
	return s.technologyRepo.Delete(id)
}

func (s *technologyService) GetActiveTechnologies() ([]models.TechnologyResponse, error) {
	technologies, err := s.technologyRepo.GetActive()
	if err != nil {
		return nil, err
	}

	var responses []models.TechnologyResponse
	for _, technology := range technologies {
		responses = append(responses, technology.ToResponse())
	}

	return responses, nil
}

func (s *technologyService) GetTechnologiesByCategory(category string) ([]models.TechnologyResponse, error) {
	technologies, err := s.technologyRepo.GetByCategory(category)
	if err != nil {
		return nil, err
	}

	var responses []models.TechnologyResponse
	for _, technology := range technologies {
		responses = append(responses, technology.ToResponse())
	}

	return responses, nil
}

func (s *technologyService) GetTechnologiesCount() (int64, error) {
	return s.technologyRepo.GetCount()
}

func (s *technologyService) UpdateTechnologiesOrder(items []OrderItem) error {
	for _, item := range items {
		if err := s.technologyRepo.UpdateOrder(item.ID, item.Order); err != nil {
			return err
		}
	}
	return nil
}
