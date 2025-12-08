package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Username           string    `json:"username" gorm:"unique;not null"`
	Email              string    `json:"email" gorm:"unique;not null"`
	Password           string    `json:"-" gorm:"not null"`        // Hide password in JSON responses
	Role               string    `json:"role" gorm:"default:user"` // Keep for backward compatibility
	RoleID             *uint     `json:"role_id" gorm:"index"`
	UserRole           *Role     `json:"user_role,omitempty" gorm:"foreignKey:RoleID"`
	IsActive           bool      `json:"is_active" gorm:"default:true"`
	MustChangePassword bool      `json:"must_change_password" gorm:"default:false"` // Force password change on first login
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token              string `json:"token"`
	User               User   `json:"user"`
	ExpiresAt          int64  `json:"expires_at"`
	MustChangePassword bool   `json:"must_change_password"` // Indicates if user must change password
}

// ChangePasswordRequest for changing user password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// HashPassword hashes the password using bcrypt
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// BeforeCreate hook to hash password before saving
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		return u.HashPassword(u.Password)
	}
	return nil
}
