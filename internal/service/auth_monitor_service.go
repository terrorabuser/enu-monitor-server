package service

import (
	"errors"
	"golang_gpt/internal/entity"
	"golang_gpt/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKeyMonitor       = []byte("your_secret_key")
	refreshTokenKey     = []byte("your_refresh_secret_key")
	accessTokenValidity = 1 * time.Hour  // Срок действия access токена - 1 час
	refreshTokenValidity = 30 * 24 * time.Hour // Срок действия refresh токена - 7 дней
)

type AuthMonitorService struct {
	monitorRepo *repository.MonitorRepository
}

func NewAuthMonitorService(monitorRepo *repository.MonitorRepository) *AuthMonitorService {
	return &AuthMonitorService{monitorRepo: monitorRepo}
}

type MonitorClaims struct {
	MacAddress string `json:"macaddress"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	MacAddress string `json:"macaddress"`
	jwt.RegisteredClaims
}

// Аутентификация монитора
func (s *AuthMonitorService) AuthenticateMonitor(macaddress, password string) (*entity.Monitor, error) {
	monitor, err := s.monitorRepo.CheckMonitorByPassword(macaddress)
	if err != nil {
		return nil, errors.New("monitor not found")
	}
	if !CheckPasswordHash(password, monitor.PasswordHash) {
		return nil, errors.New("invalid password")
	}
	return monitor, nil
}



func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// Генерация access токена
func (s *AuthMonitorService) GenerateMonitorJWT(macaddress string) (string, error) {
	expirationTime := time.Now().Add(accessTokenValidity)
	claims := &MonitorClaims{
		MacAddress: macaddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKeyMonitor)
}

// Генерация refresh токена
func (s *AuthMonitorService) GenerateRefreshToken(macaddress string) (string, error) {
	expirationTime := time.Now().Add(refreshTokenValidity)
	claims := &RefreshClaims{
		MacAddress: macaddress,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenKey)
}

// Валидация access токена
func (s *AuthMonitorService) ValidateJWT(tokenString string) (*MonitorClaims, error) {
	claims := &MonitorClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKeyMonitor, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}

// Обновление access токена с использованием refresh токена
func (s *AuthMonitorService) RefreshAccessToken(refreshToken string) (string, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	// Если refresh токен действителен, создаем новый access токен
	return s.GenerateMonitorJWT(claims.MacAddress)
}
