package service

import (
	"golang_gpt/internal/entity"
	"golang_gpt/internal/repository"
	"time"
)

type ContentService struct {
	repo *repository.ContentRepository
}

func NewContentService(repo *repository.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

// Добавление контента
func (s *ContentService) AddContent(content *entity.ContentForDB) (int, error) {

	// 0. Создаем транзацкию
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return 0, err
	}

	// 1. Добавляем контент и получаем его ID
	contentID, err := s.repo.AddContent(tx, content)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 2. history content  created
	contentHistory := entity.ContentHistory{
		ContentID: contentID,
		StatusID:  entity.ContentCreated,
		CreatedAt: time.Now(),
		UserID:    content.UserID,
	}
	
	// 3. Добавляем запись в историю
	err = s.repo.AddContentHistory(tx, &contentHistory)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	// 5. Возвращаем ID контента
	return contentID, nil
}


func (s *ContentService) GetMacAddressByLocation(building string, floor int, notes string) (string, error) {
	return s.repo.GetMacAddressByLocation(building, floor, notes)
}




// content_id, status_id, user_id
func (s *ContentService) GetContents(filter *entity.ContentFilter) ([]*entity.ContentForDB, error) {

	contents, err := s.repo.GetContents(filter)
	if err != nil {
		return nil, err
	}
	return contents, nil

}


func (s *ContentService) SendContentToModeration(contentID, statusID int, userID int64) (string, error) {
	// если переданный statusID не равен entity.ContentModerated, то возвращаем ошибку
	if statusID != entity.ContentModerated {
		return "Недостаточно прав для совершения операции", nil
	}


	// старт транзакции
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return "", err
	}
	defer tx.Rollback() // Откатится, если `Commit()` не будет вызван

	// Получаем последний status_id 
	lastStatusID, err := s.repo.GetLastStatusID(tx, contentID)
	if err != nil {
		return "", err
	}

	// Проверка lastStatusID
	if lastStatusID == entity.ContentModerated{
		return "Контент уже отправлен на модерацию", nil
	}

	if lastStatusID == entity.ContentApproved {
		return "Контент уже одобрен", nil
	}

	if lastStatusID == entity.ContentCreated || lastStatusID == entity.ContentRejected {

		err := s.repo.AddContentHistory(tx, &entity.ContentHistory{
			ContentID: contentID,
			StatusID:  statusID,
			CreatedAt: time.Now(),
			UserID:    userID,
		})
		if err != nil {
			return "", err
		}

		// Фиксация транзакции
		if err := tx.Commit(); err != nil {
			return "", err
		}
		return "Контент отправлен на модерацию!!!", nil
	}

	
	

	return "ы", nil
}


func (s *ContentService) ModerateContent(req *entity.ModerateContentRequest) (bool, error) {
	// Начало транзакции
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return false, err
	}
	defer tx.Rollback() // Откат транзакции в случае ошибки


	// Получаем последний статус контента
	lastStatusID , err := s.repo.GetLastStatusID(tx, req.ContentID)
	if err != nil {
		return false, err
	}


	// Проверяем, что статус контента позволяет его модерацию
	if lastStatusID == entity.ContentModerated {
		err := s.repo.AddContentHistory(tx, &entity.ContentHistory{
			ContentID: req.ContentID,
			StatusID:  req.StatusID,
			CreatedAt: time.Now(),
			UserID:    req.UserID,
			Reason:  req.Reason,
		})
		if err != nil {
			return false, err
		}

		// Фиксируем транзакцию перед возвратом успеха
		if err = tx.Commit(); err != nil {
			return false, err
		}

		return true, nil
	}  
	
	return false, err
}