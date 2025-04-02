package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang_gpt/internal/entity"
	"log"
)

type ContentRepository struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{
		db: db,
	}
}



// Добавление контента в базу данных
func (r *ContentRepository) AddContent(tx *sql.Tx, content *entity.ContentForDB) (int, error) {
	var id int

	err := tx.QueryRow(
		"INSERT INTO content (user_id, macaddress, file_name, file_path, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		content.UserID, content.MacAddress, content.FileName, content.FilePath, content.StartTime, content.EndTime,
	).Scan(&id)
	return id, err
}


// Получение контента по мак адресу монитора
func (r *ContentRepository) GetContentByMonitor(MacAddress string) (*entity.ContentForDB, error) {
	var content entity.ContentForDB
	err := r.db.QueryRow(
		"SELECT id, user_id, macaddress, file_name, file_path, start_time, end_time FROM content WHERE macaddress = $1",
		MacAddress,
	).Scan(&content.ID, &content.UserID, &content.MacAddress, &content.FileName, &content.FilePath, &content.StartTime, &content.EndTime)
	if err != nil {
		return nil, err
	}
	return &content, nil
}


// Получение списка контента с фильтрацией
func (r *ContentRepository) GetContents(ctx context.Context, query string, args []interface{}) ([]*entity.ContentForDB, error) {
	rows , err := r.db.QueryContext(ctx,query,args...)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var contents []*entity.ContentForDB
	for rows.Next() {
		var content entity.ContentForDB
		var history entity.ContentHistory
		
		if err := rows.Scan(
			&content.ID,
			&content.UserID,
			&content.MacAddress,
			&content.FileName,
			&content.FilePath,
			&content.StartTime,
			&content.EndTime,
			&history.ID,
			&history.ContentID,
			&history.StatusID,
			&history.CreatedAt,
			&history.UserID,
			&history.Reason,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		content.LatestHistory = &history
		contents = append(contents, &content)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return contents, nil
}


// Удаление контента по его ID
func (r *ContentRepository) DeleteContent(contentID int) error {
	_, err := r.db.Exec("DELETE FROM content WHERE id = $1", contentID)
	return err
}


// Добавление истории контента
func (r *ContentRepository) AddContentHistory(tx *sql.Tx, contentHistory *entity.ContentHistory) error {
	log.Printf("Добавление истории контента: %+v", contentHistory)
	_, err := tx.Exec(
		"INSERT INTO content_history (content_id, status_id, created_at, user_id, reason) VALUES ($1, $2, $3, $4, $5)",
		contentHistory.ContentID, contentHistory.StatusID, contentHistory.CreatedAt, contentHistory.UserID, contentHistory.Reason,
	)
	return err
}


// Получение последнего статуса контента по его ID
func (r *ContentRepository) GetLastStatusID(tx *sql.Tx, content_id int) (int, error) {
	var lastStatusID int 
	query := `SELECT status_id FROM content_history WHERE content_id = $1 ORDER BY id DESC LIMIT 1`
	err := tx.QueryRow(query, content_id).Scan(&lastStatusID)
	if err != nil {
		return 0, err
	}

	return lastStatusID, err
}


// Старт транзакции
func (r *ContentRepository) BeginTransaction() (*sql.Tx, error) {
	// создаем транзакцию
	return r.db.Begin()
}
