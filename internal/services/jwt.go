package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secretKey string
	issuer    string
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (j *JWTService) GenerateToken(userID uint, username, role string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (j *JWTService) RefreshToken(tokenString string) (string, int64, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", 0, err
	}

	// Check if token is about to expire (within 1 hour)
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", 0, errors.New("token is still valid")
	}

	return j.GenerateToken(claims.UserID, claims.Username, claims.Role)
}
