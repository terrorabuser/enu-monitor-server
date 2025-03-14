package repository

import (
	"database/sql"
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

func (r *ContentRepository) AddContent(content *entity.ContentForDB) (int, error) {
	var id int

	log.Println(content.UserID)

	err := r.db.QueryRow(
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

// Получение мак адреса по зданию и этажу
func (r *ContentRepository) GetMacAddressByLocation(building, floor, notes string) (string, error) {
	var macAddress string
	err := r.db.QueryRow(
		"SELECT macaddress FROM monitors WHERE building = $1 AND floor = $2 AND notes = $3",
		building, floor, notes,
	).Scan(&macAddress)

	if err != nil {
		return "", err
	}
	return macAddress, nil
}

func (r *ContentRepository) GetContentByID(contentID int) (*entity.ContentForDB, error) {
	var content entity.ContentForDB
	err := r.db.QueryRow(
		"SELECT id, user_id, macaddress, file_name, file_path, start_time, end_time FROM content WHERE id = $1",
		contentID,
	).Scan(&content.ID, &content.UserID, &content.MacAddress, &content.FileName, &content.FilePath, &content.StartTime, &content.EndTime)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// GetPendingContent получает список контента, ожидающего модерации
// Получение списка контента, ожидающего модерации
func (r *ContentRepository) GetContentForModeration() ([]*entity.ContentForDB, error) {
	query := `SELECT id, user_id, macaddress, file_name, file_path, start_time, end_time FROM content WHERE latest_history = 1 OR latest_history = 2`
	rows, err := r.db.Query(query)
	if err != nil { 
		log.Printf("Ошибка при получении ожидающего контента: %v", err)
		return nil, err
	}
	defer rows.Close()

	var contents []*entity.ContentForDB
	for rows.Next() {
		var content entity.ContentForDB
		err := rows.Scan(&content.ID, &content.UserID, &content.MacAddress, &content.FileName,
			&content.FilePath, &content.StartTime, &content.EndTime)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			continue // Пропускаем ошибочные строки
		}
		contents = append(contents, &content)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка при итерации по строкам: %v", err)
		return nil, err
	}
	log.Printf("Получено %d строк ожидающего контента", len(contents))
	log.Printf("Получен контент: %v", contents)
	return contents, nil
}


// UpdateContentStatus обновляет статус контента
func (r *ContentRepository) UpdateContentStatus(contentID int, status bool, reason string) error {
	_, err := r.db.Exec(
		"UPDATE moderated_content SET approved = $1, reason = $2 WHERE id = $3",
		status, reason, contentID,
	)
	return err
}

// AddModeratedContent добавляет контент в таблицу moderated_content
func (r *ContentRepository) AddModeratedContent(content *entity.ContentForDB) (int, error) {
	var id int

	// Значения NULL для approved и reason
	approved := sql.NullBool{Valid: false}    // NULL
	reason := sql.NullString{Valid: false}    // NULL

	err := r.db.QueryRow(
		"INSERT INTO moderated_content (user_id, macaddress, file_name, file_path, start_time, end_time, approved, reason) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		content.UserID, content.MacAddress, content.FileName, content.FilePath, content.StartTime, content.EndTime, approved, reason,
	).Scan(&id)

	return id, err
}


func (r *ContentRepository) DeleteContent(contentID int) error {
	_, err := r.db.Exec("DELETE FROM content WHERE id = $1", contentID)
	return err
}


func (r *ContentRepository) AddContentHistory(contentHistory *entity.ContentHistory) error {
	_, err := r.db.Exec(
		"INSERT INTO content_history (content_id, status_id, created_at, user_id) VALUES ($1, $2, $3, $4)",
		contentHistory.ContentID, contentHistory.StatusID, contentHistory.CreatedAt, contentHistory.UserID,
	)
	return err
}

func (r *ContentRepository) UpdateContentLatestHistory(contentID, statusID int) error {
	// Обновляем поле latest_history в таблице content
	_, err := r.db.Exec(
		"UPDATE content SET latest_history = $1 WHERE id = $2",
		statusID, contentID,
	)
	if err != nil {
		return err // Если первая операция не удалась, выходим
	}

	// Обновляем статус последней записи в content_history
	_, err = r.db.Exec(
		"UPDATE content_history SET status_id = $1 WHERE content_id = $2",
		statusID, contentID,
	)
	
	return err
}

