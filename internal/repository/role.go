package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) Update(id uint, role *models.Role) error {
	return r.db.Model(&models.Role{}).Where("id = ?", id).Updates(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	role := &models.Role{}
	if err := r.db.First(role, roleID).Error; err != nil {
		return err
	}

	// Clear existing permissions
	if err := r.db.Model(role).Association("Permissions").Clear(); err != nil {
		return err
	}

	// Assign new permissions
	if len(permissionIDs) > 0 {
		var permissions []models.Permission
		if err := r.db.Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}
		return r.db.Model(role).Association("Permissions").Append(permissions)
	}

	return nil
}

func (r *RoleRepository) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}
