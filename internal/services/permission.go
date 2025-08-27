package services

import (
	"errors"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"

	"gorm.io/gorm"
)

type PermissionService struct {
	permissionRepo *repository.PermissionRepository
}

func NewPermissionService(permissionRepo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
	}
}

func (s *PermissionService) CreatePermission(req *models.CreatePermissionRequest) (*models.Permission, error) {
	// Check if permission name already exists
	_, err := s.permissionRepo.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("permission name already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if resource + action combination already exists
	_, err = s.permissionRepo.GetByResourceAndAction(req.Resource, req.Action)
	if err == nil {
		return nil, errors.New("permission with this resource and action already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	permission := &models.Permission{
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		IsActive:    true,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return s.permissionRepo.GetByID(permission.ID)
}

func (s *PermissionService) GetPermission(id uint) (*models.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	return s.permissionRepo.GetAll()
}

func (s *PermissionService) GetPermissionsByResource(resource string) ([]models.Permission, error) {
	return s.permissionRepo.GetByResource(resource)
}

func (s *PermissionService) UpdatePermission(id uint, req *models.UpdatePermissionRequest) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new name already exists (if name is being changed)
	if req.Name != "" && req.Name != permission.Name {
		_, err := s.permissionRepo.GetByName(req.Name)
		if err == nil {
			return nil, errors.New("permission name already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		permission.Name = req.Name
	}

	if req.Description != "" {
		permission.Description = req.Description
	}

	// Check if resource + action combination already exists (if being changed)
	if req.Resource != "" || req.Action != "" {
		newResource := req.Resource
		if newResource == "" {
			newResource = permission.Resource
		}
		newAction := req.Action
		if newAction == "" {
			newAction = permission.Action
		}

		if newResource != permission.Resource || newAction != permission.Action {
			_, err := s.permissionRepo.GetByResourceAndAction(newResource, newAction)
			if err == nil {
				return nil, errors.New("permission with this resource and action already exists")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}

		permission.Resource = newResource
		permission.Action = newAction
	}

	if req.IsActive != nil {
		permission.IsActive = *req.IsActive
	}

	if err := s.permissionRepo.Update(id, permission); err != nil {
		return nil, err
	}

	return s.permissionRepo.GetByID(id)
}

func (s *PermissionService) DeletePermission(id uint) error {
	// Check if permission exists
	_, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.permissionRepo.Delete(id)
}

// Helper function to initialize default permissions
func (s *PermissionService) InitializeDefaultPermissions() error {
	resources := []string{"users", "roles", "permissions", "projects", "technologies", "experiences", "testimonials", "contacts", "services", "uploads"}
	actions := []string{"create", "read", "update", "delete"}

	for _, resource := range resources {
		for _, action := range actions {
			permissionName := resource + ":" + action

			// Check if permission already exists
			_, err := s.permissionRepo.GetByName(permissionName)
			if err == nil {
				continue // Permission already exists
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			// Create permission
			permission := &models.Permission{
				Name:        permissionName,
				Description: "Permission to " + action + " " + resource,
				Resource:    resource,
				Action:      action,
				IsActive:    true,
			}

			if err := s.permissionRepo.Create(permission); err != nil {
				return err
			}
		}
	}

	return nil
}
