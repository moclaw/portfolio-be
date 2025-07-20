package services

import (
	"encoding/json"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type ProjectService interface {
	CreateProject(request *models.ProjectRequest) (*models.ProjectResponse, error)
	GetAllProjects() ([]models.ProjectResponse, error)
	GetProjectByID(id uint) (*models.ProjectResponse, error)
	UpdateProject(id uint, request *models.ProjectRequest) (*models.ProjectResponse, error)
	DeleteProject(id uint) error
	GetActiveProjects() ([]models.ProjectResponse, error)
	GetProjectsCount() (int64, error)
	UpdateProjectsOrder(items []OrderItem) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

func (s *projectService) CreateProject(request *models.ProjectRequest) (*models.ProjectResponse, error) {
	// Convert tags slice to JSON string
	tagsJSON, err := json.Marshal(request.Tags)
	if err != nil {
		return nil, err
	}

	project := &models.Project{
		Name:           request.Name,
		Description:    request.Description,
		Tags:           string(tagsJSON),
		Image:          request.Image,
		SourceCodeLink: request.SourceCodeLink,
		LiveDemoLink:   request.LiveDemoLink,
		Order:          request.Order,
		IsActive:       request.IsActive,
	}

	err = s.projectRepo.Create(project)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(project)
	return &response, nil
}

func (s *projectService) GetAllProjects() ([]models.ProjectResponse, error) {
	projects, err := s.projectRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.ProjectResponse
	for _, project := range projects {
		responses = append(responses, s.convertToResponse(&project))
	}

	return responses, nil
}

func (s *projectService) GetProjectByID(id uint) (*models.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(project)
	return &response, nil
}

func (s *projectService) UpdateProject(id uint, request *models.ProjectRequest) (*models.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Convert tags slice to JSON string
	tagsJSON, err := json.Marshal(request.Tags)
	if err != nil {
		return nil, err
	}

	project.Name = request.Name
	project.Description = request.Description
	project.Tags = string(tagsJSON)
	project.Image = request.Image
	project.SourceCodeLink = request.SourceCodeLink
	project.LiveDemoLink = request.LiveDemoLink
	project.Order = request.Order
	project.IsActive = request.IsActive

	err = s.projectRepo.Update(project)
	if err != nil {
		return nil, err
	}

	response := s.convertToResponse(project)
	return &response, nil
}

func (s *projectService) DeleteProject(id uint) error {
	return s.projectRepo.Delete(id)
}

func (s *projectService) GetActiveProjects() ([]models.ProjectResponse, error) {
	projects, err := s.projectRepo.GetActive()
	if err != nil {
		return nil, err
	}

	var responses []models.ProjectResponse
	for _, project := range projects {
		responses = append(responses, s.convertToResponse(&project))
	}

	return responses, nil
}

func (s *projectService) convertToResponse(project *models.Project) models.ProjectResponse {
	// Parse JSON string back to slice
	var tags []models.ProjectTag
	if project.Tags != "" {
		json.Unmarshal([]byte(project.Tags), &tags)
	}

	return models.ProjectResponse{
		ID:             project.ID,
		Name:           project.Name,
		Description:    project.Description,
		Tags:           tags,
		Image:          project.Image,
		SourceCodeLink: project.SourceCodeLink,
		LiveDemoLink:   project.LiveDemoLink,
		Order:          project.Order,
		IsActive:       project.IsActive,
		CreatedAt:      project.CreatedAt,
		UpdatedAt:      project.UpdatedAt,
	}
}

func (s *projectService) GetProjectsCount() (int64, error) {
	return s.projectRepo.GetCount()
}

func (s *projectService) UpdateProjectsOrder(items []OrderItem) error {
	for _, item := range items {
		project, err := s.projectRepo.GetByID(item.ID)
		if err != nil {
			return err
		}

		project.Order = item.Order
		err = s.projectRepo.Update(project)
		if err != nil {
			return err
		}
	}
	return nil
}
