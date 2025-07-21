package services

import (
	"fmt"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"
	"time"
)

type ResourceService struct {
	repo       *repository.ResourceRepository
	uploadRepo *repository.UploadRepository
	s3Service  *S3Service
}

func NewResourceService(repo *repository.ResourceRepository, uploadRepo *repository.UploadRepository, s3Service *S3Service) *ResourceService {
	return &ResourceService{
		repo:       repo,
		uploadRepo: uploadRepo,
		s3Service:  s3Service,
	}
}

func (s *ResourceService) CreateResource(req *models.ResourceCreateRequest) (*models.ResourceResponse, error) {
	// Verify upload exists
	upload, err := s.uploadRepo.GetByID(req.UploadID)
	if err != nil {
		return nil, fmt.Errorf("upload not found: %w", err)
	}

	resource := &models.Resource{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Tags:        req.Tags,
		UploadID:    req.UploadID,
		Alt:         req.Alt,
		IsPublic:    true,
		IsActive:    true,
	}

	if req.IsPublic != nil {
		resource.IsPublic = *req.IsPublic
	}
	if req.IsActive != nil {
		resource.IsActive = *req.IsActive
	}

	if err := s.repo.Create(resource); err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Load the upload relationship
	resource.Upload = *upload
	response := resource.ToResponse()
	return &response, nil
}

func (s *ResourceService) GetResourceByID(id uint) (*models.ResourceResponse, error) {
	resource, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	// Increment view count
	go s.repo.IncrementViewCount(id)

	response := resource.ToResponse()
	return &response, nil
}

func (s *ResourceService) GetAllResources(limit, offset int) ([]models.ResourceResponse, error) {
	resources, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	responses := make([]models.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = resource.ToResponse()
	}

	return responses, nil
}

func (s *ResourceService) GetResourcesByType(resourceType models.ResourceType, limit, offset int) ([]models.ResourceResponse, error) {
	resources, err := s.repo.GetByType(resourceType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources by type: %w", err)
	}

	responses := make([]models.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = resource.ToResponse()
	}

	return responses, nil
}

func (s *ResourceService) GetResourcesByCategory(category string, limit, offset int) ([]models.ResourceResponse, error) {
	resources, err := s.repo.GetByCategory(category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources by category: %w", err)
	}

	responses := make([]models.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = resource.ToResponse()
	}

	return responses, nil
}

func (s *ResourceService) GetPublicResources(limit, offset int) ([]models.ResourceResponse, error) {
	resources, err := s.repo.GetPublic(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public resources: %w", err)
	}

	responses := make([]models.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = resource.ToResponse()
	}

	return responses, nil
}

func (s *ResourceService) UpdateResource(id uint, req *models.ResourceUpdateRequest) (*models.ResourceResponse, error) {
	// Check if resource exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("resource not found: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Tags != nil {
		updates["tags"] = *req.Tags
	}
	if req.Alt != nil {
		updates["alt"] = *req.Alt
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.repo.Update(id, updates); err != nil {
			return nil, fmt.Errorf("failed to update resource: %w", err)
		}
	}

	// Return updated resource
	updatedResource, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated resource: %w", err)
	}

	response := updatedResource.ToResponse()
	return &response, nil
}

func (s *ResourceService) DeleteResource(id uint) error {
	// Check if resource exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("resource not found: %w", err)
	}

	// Delete resource (upload will be handled separately if needed)
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return nil
}

func (s *ResourceService) SearchResources(query string, limit, offset int) ([]models.ResourceResponse, error) {
	resources, err := s.repo.Search(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search resources: %w", err)
	}

	responses := make([]models.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = resource.ToResponse()
	}

	return responses, nil
}

func (s *ResourceService) IncrementDownloadCount(id uint) error {
	return s.repo.IncrementDownloadCount(id)
}

// GetResourceStats returns statistics about resources
func (s *ResourceService) GetResourceStats() (map[string]interface{}, error) {
	total, err := s.repo.Count()
	if err != nil {
		return nil, err
	}

	imageCount, _ := s.repo.CountByType(models.ResourceTypeImage)
	documentCount, _ := s.repo.CountByType(models.ResourceTypeDocument)
	videoCount, _ := s.repo.CountByType(models.ResourceTypeVideo)
	audioCount, _ := s.repo.CountByType(models.ResourceTypeAudio)
	otherCount, _ := s.repo.CountByType(models.ResourceTypeOther)

	return map[string]interface{}{
		"total": total,
		"by_type": map[string]int64{
			"image":    imageCount,
			"document": documentCount,
			"video":    videoCount,
			"audio":    audioCount,
			"other":    otherCount,
		},
	}, nil
}

// RefreshExpiredURLs checks for expired uploads and refreshes their URLs
func (s *ResourceService) RefreshExpiredURLs() error {
	// Get resources with uploads expiring within 24 hours
	resources, err := s.repo.GetExpiringSoon(24 * time.Hour)
	if err != nil {
		return fmt.Errorf("failed to get expiring resources: %w", err)
	}

	for _, resource := range resources {
		// Generate new presigned URL with extended expiry
		newURL, err := s.s3Service.GeneratePresignedURL(resource.Upload.S3Key, 7*24*time.Hour) // 7 days
		if err != nil {
			fmt.Printf("Failed to generate new URL for upload %d: %v\n", resource.Upload.ID, err)
			continue
		}

		// Update upload URL and expiry
		newExpiry := time.Now().Add(7 * 24 * time.Hour)
		if err := s.uploadRepo.UpdateURL(resource.Upload.ID, newURL); err != nil {
			fmt.Printf("Failed to update URL for upload %d: %v\n", resource.Upload.ID, err)
			continue
		}

		if err := s.uploadRepo.UpdateExpiry(resource.Upload.ID, &newExpiry); err != nil {
			fmt.Printf("Failed to update expiry for upload %d: %v\n", resource.Upload.ID, err)
			continue
		}

		fmt.Printf("Refreshed URL for upload %d (resource: %s)\n", resource.Upload.ID, resource.Name)
	}

	return nil
}

// GetResourceDownloadURL generates a download URL for a resource and increments download count
func (s *ResourceService) GetResourceDownloadURL(id uint) (string, error) {
	resource, err := s.repo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("resource not found: %w", err)
	}

	// Increment download count
	go s.repo.IncrementDownloadCount(id)

	// Generate presigned URL for download (valid for 1 hour)
	downloadURL, err := s.s3Service.GeneratePresignedURL(resource.Upload.S3Key, 1*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return downloadURL, nil
}

// CountResources returns total count of resources
func (s *ResourceService) CountResources() (int64, error) {
	return s.repo.Count()
}

// CountResourcesByType returns count of resources by type
func (s *ResourceService) CountResourcesByType(resourceType models.ResourceType) (int64, error) {
	return s.repo.CountByType(resourceType)
}

// CountResourcesByCategory returns count of resources by category
func (s *ResourceService) CountResourcesByCategory(category string) (int64, error) {
	return s.repo.CountByCategory(category)
}

// CountSearchResults returns count of search results (simplified implementation)
func (s *ResourceService) CountSearchResults(query string) (int64, error) {
	// For simplicity, we'll get all search results and count them
	// In production, you might want to optimize this with a separate count query
	resources, err := s.repo.Search(query, 1000, 0) // Get up to 1000 results
	if err != nil {
		return 0, err
	}
	return int64(len(resources)), nil
}
