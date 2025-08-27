package repository

import (
	"portfolio-be/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserRole.Permissions").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserRole.Permissions").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("UserRole.Permissions").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("UserRole.Permissions").Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *UserRepository) ToggleActiveStatus(id uint, isActive bool) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *UserRepository) AssignRole(userID uint, roleID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role_id", roleID).Error
}

func (r *UserRepository) GetUserPermissions(userID uint) ([]models.Permission, error) {
	var user models.User
	err := r.db.Preload("UserRole.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	if user.UserRole == nil {
		return []models.Permission{}, nil
	}

	return user.UserRole.Permissions, nil
}

func (r *UserRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}
