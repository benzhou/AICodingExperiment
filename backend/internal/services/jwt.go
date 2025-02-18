package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenInfo struct {
	ExpiresIn int64  `json:"expires_in"` // seconds until expiration
	Token     string `json:"token"`
}

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) GenerateToken(user *models.User) (*TokenInfo, error) {
	expirationTime := time.Now().Add(time.Minute * time.Duration(config.JWTExpiryMinutes))

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Token:     tokenString,
		ExpiresIn: expirationTime.Unix() - time.Now().Unix(),
	}, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (s *JWTService) GetTokenInfo(tokenString string) (*TokenInfo, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	exp := int64((*claims)["exp"].(float64))
	return &TokenInfo{
		Token:     tokenString,
		ExpiresIn: exp - time.Now().Unix(),
	}, nil
}
