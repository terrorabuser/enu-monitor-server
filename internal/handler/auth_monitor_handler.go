package handler

import (
	"golang_gpt/internal/entity"
	"golang_gpt/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)



type AuthMonitorHandler struct {
	authMonitorService *service.AuthMonitorService
}

func NewAuthMonitorHandler(authMonitorService *service.AuthMonitorService) *AuthMonitorHandler {
	return &AuthMonitorHandler{authMonitorService: authMonitorService}
}

func (h *AuthMonitorHandler) MonitorLogin(c *gin.Context) {
	var MonitorRequest struct {
		MacAddress string `json:"macaddress"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&MonitorRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверяем логин и пароль
	// monitor, err := h.authMonitorService.AuthenticateMonitor(MonitorRequest.MacAddress, MonitorRequest.Password)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	// 	return
	// }

	monitor := &entity.Monitor{
		MacAddress: MonitorRequest.MacAddress,
	}

	if MonitorRequest.Password != "123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	

	// Генерируем access и refresh токены
	accessToken, err := h.authMonitorService.GenerateMonitorJWT(monitor.MacAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken, err := h.authMonitorService.GenerateRefreshToken(monitor.MacAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Middleware для проверки JWT и роли
func (h *AuthMonitorHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Убираем "Bearer " из начала токена
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims, err := h.authMonitorService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Логирование данных токена
		log.Printf("Authenticated Monitor: %s", claims.MacAddress)

		// Передаём MacAddress в контекст для дальнейшего использования
		c.Set("macaddress", claims.MacAddress)

		c.Next()
	}
}

func (h *AuthMonitorHandler) RefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверяем refresh токен
	newAccessToken, err := h.authMonitorService.RefreshAccessToken(refreshRequest.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Возвращаем новый access токен
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
