package handler

import (
	"fmt"
	"golang_gpt/internal/entity"
	"golang_gpt/internal/repository"
	"golang_gpt/internal/service"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

type SocketHandler struct {
	socketRepo *repository.SocketRepository
	authMonitorService *service.AuthMonitorService
}

// Конструктор
func NewSocketHandler(socketRepo *repository.SocketRepository, authMonitorService *service.AuthMonitorService) *SocketHandler {
	return &SocketHandler{socketRepo: socketRepo, authMonitorService: authMonitorService}
}

// HandleConnection обрабатывает подключение клиента
func (h *SocketHandler) HandleConnection(c socketio.Conn, clients map[string]socketio.Conn) error {
	// Получаем заголовки клиента
	headers := c.RemoteHeader()
	token := headers.Get("Authorization")
	macAddress := headers.Get("MacAddress")

	// Проверяем, есть ли необходимые данные
	if macAddress == "" {
		log.Println("Client disconnected: missing token or MAC address")
		c.Emit("error", "Missing Authorization token or MAC address")
		c.Close() // Закрываем соединение сразу
		return fmt.Errorf("missing token or MAC address")
	}

	// Проверяем токен
	claims, err := h.authMonitorService.ValidateJWT(token)
	if err != nil {
		log.Println("Client disconnected: invalid token")
		c.Emit("error", "Token is not valid")
		c.Close()
		return fmt.Errorf("invalid token")
	}

	log.Printf("Claims: %v", claims)

	// Проверяем, что это токен **монитора**
	if claims.MacAddress != macAddress {
		log.Println("Client disconnected: invalid role")
		c.Emit("error", "Access denied: only monitors are allowed")
		c.Close()
		return fmt.Errorf("access denied")
	}


	// Если токен валидный, добавляем клиента в мапу
	log.Println("Client authenticated:", macAddress)
	clients[macAddress] = c
	c.SetContext(macAddress)

	//Обновляем статус клиента в БД 
	err = h.socketRepo.SetMonitorStatus(macAddress, true)
	if err != nil {
		log.Println("Error updating monitor status:", err)
	}
	

	return nil
}

// HandleDisconnect обрабатывает отключение клиента
func (h *SocketHandler) HandleDisconnect(c socketio.Conn, reason string) {
	// Получаем MacAddress из контекста
	macAddress, ok := c.Context().(string)
	if !ok {
		log.Println("Error: MacAddress not found in context")
		return
	}

	// Ставим статус false (выключение)
	err := h.socketRepo.SetMonitorStatus(macAddress, false)
	if err != nil {
		log.Println("Error updating monitor status:", err)
	}

	log.Println("Client disconnected:", c.ID(), "Reason:", reason, "MacAddress:", macAddress)
}


// HandleGetData получает данные для клиента по MAC-адресу
func (h *SocketHandler) HandleGetData(macAddress string) ([]entity.ContentForMonitor, error) {
	// Получаем данные для клиента
	data, err := h.socketRepo.GetInfoByMac(macAddress)
	if err != nil {
		log.Println("Error fetching data:", err)
		return nil, err
	}

	// Возвращаем данные
	return data, nil
}
