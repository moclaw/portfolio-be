package services

import (
	"errors"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"

	"gorm.io/gorm"
)

type RoleService struct {
	roleRepo       *repository.RoleRepository
	permissionRepo *repository.PermissionRepository
	userRepo       *repository.UserRepository
}

func NewRoleService(roleRepo *repository.RoleRepository, permissionRepo *repository.PermissionRepository, userRepo *repository.UserRepository) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (s *RoleService) CreateRole(req *models.CreateRoleRequest) (*models.Role, error) {
	// Check if role name already exists
	_, err := s.roleRepo.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("role name already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		if err := s.roleRepo.AssignPermissions(role.ID, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return s.roleRepo.GetByID(role.ID)
}

func (s *RoleService) GetRole(id uint) (*models.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.roleRepo.GetAll()
}

func (s *RoleService) UpdateRole(id uint, req *models.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new name already exists (if name is being changed)
	if req.Name != "" && req.Name != role.Name {
		_, err := s.roleRepo.GetByName(req.Name)
		if err == nil {
			return nil, errors.New("role name already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	if err := s.roleRepo.Update(id, role); err != nil {
		return nil, err
	}

	// Update permissions if provided
	if req.PermissionIDs != nil {
		if err := s.roleRepo.AssignPermissions(id, req.PermissionIDs); err != nil {
			return nil, err
		}
	}

	return s.roleRepo.GetByID(id)
}

func (s *RoleService) DeleteRole(id uint) error {
	// Check if role exists
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if role is assigned to any users
	users, err := s.userRepo.GetAll()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.RoleID != nil && *user.RoleID == id {
			return errors.New("cannot delete role that is assigned to users")
		}
	}

	return s.roleRepo.Delete(id)
}

func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// Check if role exists
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}

	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

func (s *RoleService) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	return s.roleRepo.GetRolePermissions(roleID)
}
