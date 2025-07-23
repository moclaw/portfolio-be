package services

import (
	"encoding/json"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type ExperienceService interface {
	CreateExperience(request *models.ExperienceRequest) (*models.ExperienceResponse, error)
	GetAllExperiences() ([]models.ExperienceResponse, error)
	GetExperienceByID(id uint) (*models.ExperienceResponse, error)
	UpdateExperience(id uint, request *models.ExperienceRequest) (*models.ExperienceResponse, error)
	DeleteExperience(id uint) error
	GetActiveExperiences() ([]models.ExperienceResponse, error)
	GetExperiencesCount() (int64, error)
	UpdateExperiencesOrder(items []OrderItem) error
}

type experienceService struct {
	experienceRepo repository.ExperienceRepository
}

func NewExperienceService(experienceRepo repository.ExperienceRepository) ExperienceService {
	return &experienceService{experienceRepo: experienceRepo}
}

func (s *experienceService) CreateExperience(request *models.ExperienceRequest) (*models.ExperienceResponse, error) {
	// Convert points slice to JSON string
	pointsJSON, err := json.Marshal(request.Points)
	if err != nil {
		return nil, err
	}

	experience := &models.Experience{
		Title:       request.Title,
		CompanyName: request.CompanyName,
		Icon:        request.Icon,
		IconBg:      request.IconBg,
		Date:        request.Date,
		Points:      string(pointsJSON),
		Order:       request.Order,
		IsActive:    request.IsActive,
	}

	err = s.experienceRepo.Create(experience)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(experience)
	return &response, nil
}

func (s *experienceService) GetAllExperiences() ([]models.ExperienceResponse, error) {
	experiences, err := s.experienceRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.ExperienceResponse
	for _, experience := range experiences {
		responses = append(responses, s.convertToResponse(&experience))
	}

	return responses, nil
}

func (s *experienceService) GetExperienceByID(id uint) (*models.ExperienceResponse, error) {
	experience, err := s.experienceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(experience)
	return &response, nil
}

func (s *experienceService) UpdateExperience(id uint, request *models.ExperienceRequest) (*models.ExperienceResponse, error) {
	experience, err := s.experienceRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Convert points slice to JSON string
	pointsJSON, err := json.Marshal(request.Points)
	if err != nil {
		return nil, err
	}

	experience.Title = request.Title
	experience.CompanyName = request.CompanyName
	experience.Icon = request.Icon
	experience.IconBg = request.IconBg
	experience.Date = request.Date
	experience.Points = string(pointsJSON)
	experience.Order = request.Order
	experience.IsActive = request.IsActive

	err = s.experienceRepo.Update(experience)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(experience)
	return &response, nil
}

func (s *experienceService) DeleteExperience(id uint) error {
	return s.experienceRepo.Delete(id)
}

func (s *experienceService) GetActiveExperiences() ([]models.ExperienceResponse, error) {
	experiences, err := s.experienceRepo.GetActive()
	if err != nil {
		return nil, err
	}

	var responses []models.ExperienceResponse
	for _, experience := range experiences {
		responses = append(responses, s.convertToResponse(&experience))
	}

	return responses, nil
}

func (s *experienceService) convertToResponse(experience *models.Experience) models.ExperienceResponse {
	// Parse JSON string back to slice
	var points []string
	if experience.Points != "" {
		json.Unmarshal([]byte(experience.Points), &points)
	}

	return models.ExperienceResponse{
		ID:          experience.ID,
		Title:       experience.Title,
		CompanyName: experience.CompanyName,
		Icon:        experience.Icon,
		IconBg:      experience.IconBg,
		Date:        experience.Date,
		Points:      points,
		Order:       experience.Order,
		IsActive:    experience.IsActive,
		CreatedAt:   experience.CreatedAt,
		UpdatedAt:   experience.UpdatedAt,
	}
}

func (s *experienceService) GetExperiencesCount() (int64, error) {
	return s.experienceRepo.GetCount()
}

func (s *experienceService) UpdateExperiencesOrder(items []OrderItem) error {
	for _, item := range items {
		experience, err := s.experienceRepo.GetByID(item.ID)
		if err != nil {
			return err
		}

		experience.Order = item.Order
		err = s.experienceRepo.Update(experience)
		if err != nil {
			return err
		}
	}
	return nil
}
