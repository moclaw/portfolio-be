package services

import (
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type ServiceService interface {
	CreateService(request *models.ServiceRequest) (*models.ServiceResponse, error)
	GetAllServices() ([]models.ServiceResponse, error)
	GetServiceByID(id uint) (*models.ServiceResponse, error)
	UpdateService(id uint, request *models.ServiceRequest) (*models.ServiceResponse, error)
	DeleteService(id uint) error
	GetActiveServices() ([]models.ServiceResponse, error)
	GetServicesCount() (int64, error)
	UpdateServicesOrder(items []OrderItem) error
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository) ServiceService {
	return &serviceService{serviceRepo: serviceRepo}
}

func (s *serviceService) CreateService(request *models.ServiceRequest) (*models.ServiceResponse, error) {
	service := &models.Service{
		Title:    request.Title,
		Icon:     request.Icon,
		Order:    request.Order,
		IsActive: request.IsActive,
	}

	err := s.serviceRepo.Create(service)
	if err != nil {
		return nil, err
	}

	response := service.ToResponse()
	return &response, nil
}

func (s *serviceService) GetAllServices() ([]models.ServiceResponse, error) {
	services, err := s.serviceRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.ServiceResponse
	for _, service := range services {
		responses = append(responses, service.ToResponse())
	}

	return responses, nil
}

func (s *serviceService) GetServiceByID(id uint) (*models.ServiceResponse, error) {
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := service.ToResponse()
	return &response, nil
}

func (s *serviceService) UpdateService(id uint, request *models.ServiceRequest) (*models.ServiceResponse, error) {
	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	service.Title = request.Title
	service.Icon = request.Icon
	service.Order = request.Order
	service.IsActive = request.IsActive

	err = s.serviceRepo.Update(service)
	if err != nil {
		return nil, err
	}

	response := service.ToResponse()
	return &response, nil
}

func (s *serviceService) DeleteService(id uint) error {
	return s.serviceRepo.Delete(id)
}

func (s *serviceService) GetActiveServices() ([]models.ServiceResponse, error) {
	services, err := s.serviceRepo.GetActive()
	if err != nil {
		return nil, err
	}

	var responses []models.ServiceResponse
	for _, service := range services {
		responses = append(responses, service.ToResponse())
	}

	return responses, nil
}

func (s *serviceService) GetServicesCount() (int64, error) {
	return s.serviceRepo.GetCount()
}

func (s *serviceService) UpdateServicesOrder(items []OrderItem) error {
	for _, item := range items {
		service, err := s.serviceRepo.GetByID(item.ID)
		if err != nil {
			return err
		}

		service.Order = item.Order
		err = s.serviceRepo.Update(service)
		if err != nil {
			return err
		}
	}
	return nil
}
