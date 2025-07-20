package services

import (
	"fmt"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
)

type ContentService struct {
	repo *repository.ContentRepository
}

func NewContentService(repo *repository.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

func (s *ContentService) CreateContent(req models.ContentRequest) (*models.ContentResponse, error) {
	content := &models.Content{
		Title:       req.Title,
		Description: req.Description,
		Body:        req.Body,
		Category:    req.Category,
		Tags:        req.Tags,
		Status:      req.Status,
		ImageURL:    req.ImageURL,
	}

	// Set default status if not provided
	if content.Status == "" {
		content.Status = "draft"
	}

	if err := s.repo.Create(content); err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	response := content.ToResponse()
	return &response, nil
}

func (s *ContentService) GetContentByID(id uint) (*models.ContentResponse, error) {
	content, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	response := content.ToResponse()
	return &response, nil
}

func (s *ContentService) GetAllContent(limit, offset int, category, status string) ([]models.ContentResponse, error) {
	var contents []models.Content
	var err error

	if category != "" {
		contents, err = s.repo.GetByCategory(category, limit, offset)
	} else if status != "" {
		contents, err = s.repo.GetByStatus(status, limit, offset)
	} else {
		contents, err = s.repo.GetAll(limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get contents: %w", err)
	}

	responses := make([]models.ContentResponse, len(contents))
	for i, content := range contents {
		responses[i] = content.ToResponse()
	}

	return responses, nil
}

func (s *ContentService) UpdateContent(id uint, req models.ContentRequest) (*models.ContentResponse, error) {
	content, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// Update fields
	content.Title = req.Title
	content.Description = req.Description
	content.Body = req.Body
	content.Category = req.Category
	content.Tags = req.Tags
	content.Status = req.Status
	content.ImageURL = req.ImageURL

	if err := s.repo.Update(content); err != nil {
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	response := content.ToResponse()
	return &response, nil
}

func (s *ContentService) DeleteContent(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete content: %w", err)
	}

	return nil
}

func (s *ContentService) GetContentCount() (int64, error) {
	return s.repo.Count()
}