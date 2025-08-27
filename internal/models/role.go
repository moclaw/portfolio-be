package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Many-to-many relationship with permissions
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`

	// One-to-many relationship with users
	Users []User `json:"users,omitempty" gorm:"foreignKey:RoleID"`
}

type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"` // e.g., "users", "projects", "technologies"
	Action      string    `json:"action"`   // e.g., "create", "read", "update", "delete"
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Many-to-many relationship with roles
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

type RolePermission struct {
	RoleID       uint      `json:"role_id" gorm:"primaryKey"`
	PermissionID uint      `json:"permission_id" gorm:"primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
}

// Request/Response DTOs
type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids"`
	IsActive      *bool  `json:"is_active"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	IsActive    *bool  `json:"is_active"`
}

type UserRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

// Helper methods
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.Name == "" {
		return errors.New("role name is required")
	}
	return nil
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.Name == "" {
		return errors.New("permission name is required")
	}
	if p.Resource == "" {
		return errors.New("permission resource is required")
	}
	if p.Action == "" {
		return errors.New("permission action is required")
	}
	return nil
}
