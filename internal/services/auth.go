package services

import (
	"errors"
	"portfolio-be/internal/models"
	"portfolio-be/internal/repository"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if email already exists
	_, err = s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, expiresAt, err := s.jwtService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) RefreshToken(tokenString string) (*models.LoginResponse, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Get user to make sure they still exist and are active
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate new token
	token, expiresAt, err := s.jwtService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}
