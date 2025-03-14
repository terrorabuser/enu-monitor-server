package server

import (
	"fmt"
	"golang_gpt/internal/handler"
	"log"

	"github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// Мапа для отслеживания подключений
var clients = make(map[string]socketio.Conn)




func NewSocketServer(socketHandler *handler.SocketHandler) *socketio.Server {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{},   // Поддержка polling
			&websocket.Transport{}, // Поддержка WebSocket
		},
	})

	// Обработчик подключения, передаем мапу для обновления
	server.OnConnect("/", func(s socketio.Conn) error {
		return socketHandler.HandleConnection(s, clients) // Передаем мапу clients
	})


	server.OnEvent("/", "joinRoom", func(s socketio.Conn, macAddress string) {
		// Получаем MAC-адрес, который был передан в заголовке при подключении
		storedMac, ok := s.Context().(string)
		if !ok || storedMac != macAddress {
			log.Println("Ошибка: MAC-адрес из сообщения не совпадает с заголовком")
			s.Emit("error", "Invalid MacAddress")
			return
		}
	
		s.Join(macAddress)
		fmt.Println("Пользователь", s.ID(), "вошёл в комнату", macAddress)
	
		// После подключения клиента, отправляем данные в комнату
		SendDataToRoom(server, socketHandler, macAddress)
	})

	server.OnDisconnect("/", socketHandler.HandleDisconnect)

	// Запуск сервера
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO error: %v", err)
		}
	}()

	return server
}

func SendDataToRoom(server *socketio.Server, socketHandler *handler.SocketHandler, macAddress string) {
	// Получаем данные для клиента по MAC-адресу
	data, err := socketHandler.HandleGetData(macAddress)
	if err != nil {
		log.Println("Error fetching data:", err)
		return
	}

	// Отправляем данные всем клиентам в комнате с именем macAddress
	server.BroadcastToRoom("/", macAddress, "data", data)
}
