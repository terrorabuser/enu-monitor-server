package repository

import (
	"database/sql"
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
func (r *ContentRepository) GetContents(filter *entity.ContentFilter) ([]*entity.ContentForDB, error) {
	query := `SELECT 
		c.id, 
		c.user_id, 
		c.macaddress, 
		c.file_name, 
		c.file_path, 
		c.start_time, 
		c.end_time,
		ch.id AS history_id,
		ch.status_id, 
		ch.created_at, 
		ch.user_id AS history_user_id
	FROM 
		content c
	LEFT JOIN 
		content_history ch 
		ON ch.content_id = c.id
		AND ch.id = (
			SELECT MAX(id) 
			FROM content_history 
			WHERE content_id = c.id
		)
	WHERE 1=1;`

	var args []interface{}

	//1, 2, 4,      2, 4, 2, 3

	log.Printf("Фильтр: %+v", filter)

	// Фильтр по UserId
	if filter.UserId != nil && *filter.UserId != 0 {
		query += fmt.Sprintf(" AND c.user_id = $%d", len(args)+1)
		args = append(args, *filter.UserId)
		log.Printf("Добавлен фильтр по UserId: %v", *filter.UserId)
	}

	// Фильтр по StatusId (одно значение)
	if filter.StatusId != nil && *filter.StatusId != 0 {
		query += fmt.Sprintf(" AND ch.status_id = $%d", len(args)+1)
		args = append(args, *filter.StatusId)
		log.Printf("Добавлен фильтр по StatusId: %v", *filter.StatusId)
	}

	// Фильтр по StartTime
	if filter.StartTime != nil && *filter.StartTime != "" {
		query += fmt.Sprintf(" AND c.start_time >= $%d", len(args)+1)
		args = append(args, *filter.StartTime)
		log.Printf("Добавлен фильтр по StartTime: %v", *filter.StartTime)
	}

	// Фильтр по EndTime
	if filter.EndTime != nil && *filter.EndTime != "" {
		query += fmt.Sprintf(" AND c.end_time <= $%d", len(args)+1)
		args = append(args, *filter.EndTime)
		log.Printf("Добавлен фильтр по EndTime: %v", *filter.EndTime)
	}

	// Сортировка по ID
	query += " ORDER BY c.id ASC"

	log.Printf("Сформированный запрос: %s", query)
	log.Printf("Аргументы для запроса: %v", args)

	// Выполнение запроса
	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("Ошибка при получении контента: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Чтение результатов
	var contents []*entity.ContentForDB
	for rows.Next() {
		var content entity.ContentForDB
		var history entity.ContentHistory
		err := rows.Scan(&content.ID, &content.UserID, &content.MacAddress, &content.FileName,
			&content.FilePath, &content.StartTime, &content.EndTime,
			&history.ID, &history.StatusID, &history.CreatedAt, &history.UserID)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			continue // Пропускаем ошибочные строки
		}


		// Присваиваем LatestHistory
		content.LatestHistory = entity.ContentHistory{
			ID:        history.ID,
			ContentID: content.ID,
			StatusID:  history.StatusID,
			CreatedAt: history.CreatedAt,
			UserID:    history.UserID,
		}

		contents = append(contents, &content)
	}

	// Проверка ошибок при итерации по строкам
	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, err
	}

	log.Printf("Получено %d строк контента", len(contents))
	return contents, nil
}


// Удаление контента по его ID
func (r *ContentRepository) DeleteContent(contentID int) error {
	_, err := r.db.Exec("DELETE FROM content WHERE id = $1", contentID)
	return err
}


// Добавление истории контента
func (r *ContentRepository) AddContentHistory(tx *sql.Tx, contentHistory *entity.ContentHistory) error {
	_, err := tx.Exec(
		"INSERT INTO content_history (content_id, status_id, created_at, user_id) VALUES ($1, $2, $3, $4)",
		contentHistory.ContentID, contentHistory.StatusID, contentHistory.CreatedAt, contentHistory.UserID,
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
