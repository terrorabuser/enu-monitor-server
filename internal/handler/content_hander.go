package handler

import (
	"context"
	"golang_gpt/internal/entity"
	pb "golang_gpt/internal/proto"
	"golang_gpt/internal/service"
	"log"

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


	// Получаем мак адресс на основе здания и этажа и примечания
	macAddress, err := h.service.GetMacAddressByLocation(req.Building, int(req.Floor), req.Notes)
	if err != nil {
		log.Println("Ошибка получения MAC-адреса:", err)
		return nil, err
	}

	content := &entity.ContentForDB{
		UserID:    req.UserId, // В будущем передавать из аутентификации
		MacAddress: macAddress,
		FileName:  req.FileName,
		FilePath:  req.FilePath,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	log.Printf("Добавляем контент: %+v", content)

	// Добавляем контент в базу данных
	id, err := h.service.AddContent(content)
	if err != nil {
		return nil, err
	}

	return &pb.AddContentResponse{ContentId: int32(id)}, nil
}


// Получение контента 
func (h *ContentHandler) GetContents(ctx context.Context, req *pb.GetContentsRequest) (*pb.GetContentsResponse, error) {
	
 
	 // Создаем фильтр из запроса
	 filter := &entity.ContentFilter{
		 UserId:    &req.UserId,
		 StatusId: &req.StatusId, 
		 StartTime: &req.StartTime,
		 EndTime:   &req.EndTime,
	 }

	 contents, err := h.service.GetContents(filter)
	 if err != nil {
		 log.Println("Ошибка получения контента для модерации:", err)
		 return nil, err
	 }
	 
	 // Логируем содержимое contents
	 log.Printf("Полученные данные о контенте: %+v", contents)
	 
	 var pbContents []*pb.ContentForDB
	 for _, content := range contents {
		 log.Printf("Обрабатываем контент: %+v", content)
	 
	 
		 pbContents = append(pbContents, &pb.ContentForDB{
			 Id:        int32(content.ID),
			 UserId:    content.UserID,
			 MacAddress: content.MacAddress,
			 FileName:  content.FileName,
			 FilePath:  content.FilePath,
			 StartTime: content.StartTime,
			 EndTime:   content.EndTime,
			 LatestHistory: &pb.ContentHistory{
				 Id:        int32(content.LatestHistory.ID),
				 ContentId: int32(content.LatestHistory.ContentID),
				 StatusId:  int32(content.LatestHistory.StatusID),
				 CreatedAt: content.LatestHistory.CreatedAt.String(), // Преобразуем time.Time в строку
				 UserId:    content.LatestHistory.UserID,
			 },
		 })
	 }
	 
	 return &pb.GetContentsResponse{Contents: pbContents}, nil
	 
}


// Модерация контента
func (h *ContentHandler) ModerateContent(ctx context.Context, req *pb.ModerateContentRequest) (*pb.ModerateContentResponse, error) {
	log.Printf("Received ModerateContent request: %+v", req)

	reqContent := &entity.ModerateContentRequest{
		ContentID: int(req.ContentId),
		StatusID:  int(req.StatusId),
		UserID:   req.UserId,
		Reason:   req.Reason,
	}

	// Модерируем контент
	success, err := h.service.ModerateContent(reqContent) 
	if err != nil {
		log.Println("Ошибка модерации контента:", err)
		return nil, err
	}
	

	return &pb.ModerateContentResponse{Success: success}, nil
}


func (h *ContentHandler) SendContentToModeration(ctx context.Context, req *pb.SendContentToModerationRequest) (*pb.SendContentToModerationResponse, error) {
	
	log.Printf("Received SendContentToModeration request: %+v", req)
	// возвращаем ошибку если не удалось отправить контент на модерацию

	response, err := h.service.SendContentToModeration(int(req.ContentId),int(req.StatusId), req.UserId)
	if err != nil {
		log.Println("Ошибка отправки контента на модерацию:", err)
		return &pb.SendContentToModerationResponse{Message: response}, err
	}


	return &pb.SendContentToModerationResponse{Message: response}, nil
}

