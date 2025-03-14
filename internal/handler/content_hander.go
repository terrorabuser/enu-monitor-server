package handler

import (
	"context"
	"fmt"
	"log"
	"golang_gpt/internal/entity"
	pb "golang_gpt/internal/proto"
	"golang_gpt/internal/service"

	socketio "github.com/googollee/go-socket.io"
)

type ContentHandler struct {
	pb.UnimplementedContentServiceServer
	service *service.ContentService
	server  *socketio.Server
}

func NewContentHandler(svc *service.ContentService, server *socketio.Server) *ContentHandler {
	return &ContentHandler{service: svc, server: server}
}

// Добавление контента
func (h *ContentHandler) AddContent(ctx context.Context, req *pb.AddContentRequest) (*pb.AddContentResponse, error) {
	log.Printf("Received AddContent request: %+v", req)

	macAddress, err := h.service.GetMacAddressByLocation(req.Building, req.Floor, req.Notes)
	if err != nil {
		log.Println("Ошибка получения MAC-адреса:", err)
		return nil, err
	}

	content := &entity.ContentForDB{
		UserID:    1, // В будущем передавать из аутентификации
		MacAddress: macAddress,
		FileName:  req.FileName,
		FilePath:  req.FilePath,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	id, err := h.service.AddContent(content)
	if err != nil {
		return nil, err
	}

	return &pb.AddContentResponse{ContentId: fmt.Sprintf("%d", id)}, nil
}

// Получение немодерированного контента
func (h *ContentHandler) GetUnmoderatedContent(ctx context.Context, req *pb.GetUnmoderatedContentRequest) (*pb.GetUnmoderatedContentResponse, error) {
	log.Printf("Received GetContentForModeration request: %+v", req)

	contents, err := h.service.GetContentForModeration()
	if err != nil {
		log.Println("Ошибка получения контента для модерации:", err)
		return nil, err
	}

	var pbContents []*pb.ContentForDB
	for _, content := range contents {
		pbContents = append(pbContents, &pb.ContentForDB{
			Id:        int32(content.ID),
			UserId:    content.UserID,
			MacAddress: content.MacAddress,
			FileName:  content.FileName,
			FilePath:  content.FilePath,
			StartTime: content.StartTime,
			EndTime:   content.EndTime,
		})
	}

	return &pb.GetUnmoderatedContentResponse{Contents: pbContents}, nil
}

// Модерация контента
func (h *ContentHandler) ModerateContent(ctx context.Context, req *pb.ModerateContentRequest) (*pb.ModerateContentResponse, error) {
	log.Printf("Received ModerateContent request: %+v", req)


	log.Printf("ContentId: %d, StatusId: %d", req.ContentId, req.StatusId)
	err := h.service.UpdateContentLatestHistory(int(req.ContentId), int(req.StatusId))
	if err != nil {
		log.Printf("Ошибка обновления статуса контента ID %d: %v", req.ContentId, err)
		return nil, err
	}

	// Если контент одобрен - отправляем его на монитор
	if req.StatusId == entity.ContentApproved {
		content, _ := h.service.GetContentByID(int(req.ContentId))
		approvedContent := &entity.ContentForMonitor{
			FileName:  content.FileName,
			FilePath:  content.FilePath,
			StartTime: content.StartTime,
			EndTime:   content.EndTime,
		}

		// Отправляем данные на устройство
		h.server.BroadcastToRoom("/", content.MacAddress, "data", approvedContent)
	}

	// Проверяем статус контента
	success := false
	if req.StatusId == 3 {
		success = true
	} else if req.StatusId == 4 {
		success = false
	}

	return &pb.ModerateContentResponse{Success: success}, nil
}
