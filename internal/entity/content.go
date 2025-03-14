// Модель задания на отображение контента (ID, MonitorID, ContentURL, Schedule)
package entity

import "time"

// Константы статусов контента
const (
	ContentCreated   = iota + 1 // 1 - Создано
	ContentModerated            // 2 - Отправлено на проверку
	ContentApproved             // 3 - Принято
	ContentRejected             // 4 - Отклонено
)

type ContentForDB struct {
	ID            int            `json:"id"`
	UserID        int64          `json:"user_id"`
	MacAddress    string         `json:"macaddress"`
	FileName      string         `json:"file_name"`
	FilePath      string         `json:"file_path"`
	StartTime     string         `json:"start_time"`
	EndTime       string         `json:"end_time"`
	LatestHistory ContentHistory `json:"latest_history"`
}

type ContentHistory struct {
	ID        int       `json:"id"`
	ContentID int       `json:"content_id"`
	StatusID  int       `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int64       `json:"user_id"`
}


type ContentForMonitor struct {
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ContentRequest struct {
	Building  string `json:"building"`
	Floor     string `json:"floor"`
	Notes     string `json:"notes"`
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}


type ModerateContentRequest struct {
	ContentID int    `json:"content_id"` // ID контента, который модератор проверяет
	StatusID  int    `json:"status_id"`  // ID статуса (например, ContentApproved, ContentRejected)
	Reason    string `json:"reason,omitempty"` // Причина отклонения, если контент не принят
}