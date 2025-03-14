package server

import (
	"log"
	"net"
	"golang_gpt/internal/handler"
	pb "golang_gpt/internal/proto"
	"golang_gpt/internal/service"

	"google.golang.org/grpc"
	socketio "github.com/googollee/go-socket.io"
)

func RunGRPCServer(contentService *service.ContentService, server *socketio.Server) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	contentHandler := handler.NewContentHandler(contentService, server)
	pb.RegisterContentServiceServer(grpcServer, contentHandler)

	log.Println("gRPC сервер запущен на порту 50051...")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}
