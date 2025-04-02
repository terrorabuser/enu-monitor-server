package entity

import (
    "time"
)

// Константы статусов остаются без изменений
const (
    ContentCreated   = iota + 1
    ContentModerated
    ContentApproved
    ContentRejected
)

type ContentForDB struct {
    ID            int            `json:"id"`
    UserID        int64          `json:"user_id"`
    MacAddress    string         `json:"macaddress"`
    FileName      string         `json:"file_name"`
    FilePath      string         `json:"file_path"`
    StartTime     time.Time      `json:"start_time"`  // Изменено на time.Time
    EndTime       time.Time      `json:"end_time"`    // Изменено на time.Time
    LatestHistory *ContentHistory `json:"latest_history"`
}

type ContentHistory struct {
    ID        int       `json:"id"`
    ContentID int       `json:"content_id"`
    StatusID  int       `json:"status_id"`
    CreatedAt time.Time `json:"created_at"` // Уже правильный тип
    UserID    int64     `json:"user_id"`
    Reason    string    `json:"reason"`
}

type ContentForMonitor struct {
    FileName  string    `json:"file_name"`
    FilePath  string    `json:"file_path"`
    StartTime time.Time `json:"start_time"` // Изменено на time.Time
    EndTime   time.Time `json:"end_time"`   // Изменено на time.Time
}

type ContentRequest struct {
    Building  string    `json:"building"`
    Floor     string    `json:"floor"`
    Notes     string    `json:"notes"`
    FileName  string    `json:"file_name"`
    FilePath  string    `json:"file_path"`
    StartTime time.Time `json:"start_time"` // Изменено на time.Time
    EndTime   time.Time `json:"end_time"`   // Изменено на time.Time
}

type ModerateContentRequest struct {
    ContentID int       `json:"content_id"`
    StatusID  int       `json:"status_id"`
    UserID    int64     `json:"user_id"`
    Reason    string    `json:"reason,omitempty"`
    CreatedAt time.Time `json:"created_at"` // Добавлено для отслеживания времени модерации
}

type ContentFilter struct {
    UserId    *int64     `json:"user_id"`
    StatusId  *int32     `json:"status_id"`
    StartTime *time.Time `json:"start_time"` // Изменено на *time.Time
    EndTime   *time.Time `json:"end_time"`   // Изменено на *time.Time
}