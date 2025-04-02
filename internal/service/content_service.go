package service

import (
	"context"
	"errors"
	"fmt"
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
func (s *ContentService) GetContents(ctx context.Context, filter *entity.ContentFilter) ([]*entity.ContentForDB, error) {
	// Валидация параметров
	if filter.UserId != nil && *filter.UserId < 0 {
		return nil, errors.New("user_id must be positive")
	}

	if filter.StatusId != nil && (*filter.StatusId < 1 || *filter.StatusId > 4) {
		return nil, errors.New("status_id must be between 1 and 4")
	}


	// Построение запроса
	qb := repository.NewContentQueryBuilder().
		ApplyUserId(filter.UserId).
		ApplyStatusId(filter.StatusId)

	// Проверяем, что время не является "нулевым" 
	if filter.StartTime != nil || filter.EndTime != nil  {
		qb.ApplyTimeFilters(filter.StartTime, filter.EndTime)
	}

	query, args := qb.Build()

	// Выполнение запроса
	contents, err := s.repo.GetContents(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get contents: %w", err)
	}

	return contents, nil
	
}
func (s *ContentService) SendContentToModeration(contentID int, userID int64) (string, error) {
	// Проверка валидности входных данных
	if contentID <= 0 {
		return "", errors.New("contentID must be positive")
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
			StatusID:  entity.ContentModerated,
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