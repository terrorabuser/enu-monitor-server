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
        StatusID: entity.ContentCreated,
        CreatedAt: time.Now(),
        UserID: content.UserID,
    }

    // 3. Добавляем запись в историю
    err = s.repo.AddContentHistory(tx, &contentHistory)
    if err != nil {
        tx.Rollback()
        return 0, err
    }

    // 4. Обновляем latest_history в monitors
	err = s.repo.UpdateContentLatestHistory(tx, contentID, contentHistory.StatusID)
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



func (s *ContentService) AddModeratedContent(content *entity.ContentForDB) (int, error) {
    return s.repo.AddModeratedContent(content)
}

func (s *ContentService) GetMacAddressByLocation(building, floor, notes string) (string, error) {
    return s.repo.GetMacAddressByLocation(building, floor, notes)
}


// Подтверждение контента
// func (s *ContentService) ApproveContent(contentID int64) error {
//     return s.repo.UpdateContentStatus(contentID, "approved", "")
// }

// // Отклонение контента с указанием причины
// func (s *ContentService) RejectContent(contentID int64, reason string) error {
//     return s.repo.UpdateContentStatus(contentID, "rejected", reason)
// }

// Получение контента по ID
func (s *ContentService) GetContentByID(contentID int) (*entity.ContentForDB, error) {
    return s.repo.GetContentByID(contentID)
}

func (s *ContentService) GetContentForModeration() ([]*entity.ContentForDB, error) {
    tx, err := s.repo.BeginTransaction()
    if err != nil {
        return nil, err
    }



	contents, err := s.repo.GetContentForModeration()
    if err != nil {
        return nil, err
    }

    // Меняем статус на Отправлено на модерацию
    for _, content := range contents {
        err = s.repo.UpdateContentLatestHistory(tx, content.ID, entity.ContentModerated)
        if err != nil {
            tx.Rollback()
            return nil, err
        }
    }

    err = tx.Commit()
    if err != nil {
        return nil, err
    }

    return contents, nil


}

// Обновление статуса контента
func (s *ContentService) UpdateContentLatestHistory(contentID, statusID int) error {
    tx, err := s.repo.BeginTransaction()
    if err != nil {
        return err
    }

    err = s.repo.UpdateContentLatestHistory(tx, contentID, statusID)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}
