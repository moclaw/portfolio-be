package middleware

import (
	"fmt"
	"net/http"
	"portfolio-be/internal/repository"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PermissionMiddleware struct {
	userRepo *repository.UserRepository
}

func NewPermissionMiddleware(userRepo *repository.UserRepository) *PermissionMiddleware {
	return &PermissionMiddleware{
		userRepo: userRepo,
	}
}

// RequirePermission checks if the authenticated user has the required permission
func (m *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Convert userID to uint
		id, ok := userID.(uint)
		if !ok {
			// Try to convert from string if it's stored as string
			if strID, isString := userID.(string); isString {
				if parsedID, err := strconv.ParseUint(strID, 10, 32); err == nil {
					id = uint(parsedID)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
				c.Abort()
				return
			}
		}

		// Check if user has permission
		hasPermission, err := m.checkUserPermission(id, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if the authenticated user has any of the specified permissions
func (m *PermissionMiddleware) RequireAnyPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		id, ok := userID.(uint)
		if !ok {
			if strID, isString := userID.(string); isString {
				if parsedID, err := strconv.ParseUint(strID, 10, 32); err == nil {
					id = uint(parsedID)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
				c.Abort()
				return
			}
		}

		// Check if user has any of the required permissions
		for _, permission := range permissions {
			parts := strings.Split(permission, ":")
			if len(parts) != 2 {
				continue
			}

			hasPermission, err := m.checkUserPermission(id, parts[0], parts[1])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
				c.Abort()
				return
			}

			if hasPermission {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// RequireRole checks if the authenticated user has the required role
func (m *PermissionMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		id, ok := userID.(uint)
		if !ok {
			if strID, isString := userID.(string); isString {
				if parsedID, err := strconv.ParseUint(strID, 10, 32); err == nil {
					id = uint(parsedID)
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
				c.Abort()
				return
			}
		}

		user, err := m.userRepo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			c.Abort()
			return
		}

		// Check legacy role field first
		if user.Role == roleName {
			c.Next()
			return
		}

		// Check new role relationship
		if user.UserRole != nil && user.UserRole.Name == roleName {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
		c.Abort()
	}
}

// checkUserPermission checks if a user has a specific permission
func (m *PermissionMiddleware) checkUserPermission(userID uint, resource, action string) (bool, error) {
	user, err := m.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	// Admin role has all permissions (legacy support)
	if user.Role == "admin" {
		return true, nil
	}

	// Check permissions through role
	if user.UserRole != nil {
		// Admin role has all permissions
		if user.UserRole.Name == "admin" {
			return true, nil
		}

		for _, permission := range user.UserRole.Permissions {
			if permission.Resource == resource && permission.Action == action && permission.IsActive {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserPermissions returns all permissions for a user
func (m *PermissionMiddleware) GetUserPermissions(userID uint) ([]string, error) {
	user, err := m.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	var permissions []string

	// Admin role has all permissions
	if user.Role == "admin" || (user.UserRole != nil && user.UserRole.Name == "admin") {
		// Return all possible permissions for admin
		resources := []string{"users", "roles", "permissions", "projects", "technologies", "experiences", "testimonials", "contacts", "services", "uploads"}
		actions := []string{"create", "read", "update", "delete"}

		for _, resource := range resources {
			for _, action := range actions {
				permissions = append(permissions, fmt.Sprintf("%s:%s", resource, action))
			}
		}
		return permissions, nil
	}

	// Get permissions from role
	if user.UserRole != nil {
		for _, permission := range user.UserRole.Permissions {
			if permission.IsActive {
				permissions = append(permissions, fmt.Sprintf("%s:%s", permission.Resource, permission.Action))
			}
		}
	}

	return permissions, nil
}
