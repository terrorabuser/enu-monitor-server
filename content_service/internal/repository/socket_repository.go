package repository

import (
	// "golang_gpt/internal/entity"
	"database/sql"
	"golang_gpt/internal/entity"
)

type SocketRepository struct {
	db *sql.DB
}

func NewSocketRepository(db *sql.DB) *SocketRepository {
	return &SocketRepository{db: db}
}


func (r *SocketRepository) GetInfoByMac(MacAddress string) ([]entity.ContentForMonitor, error) {
	var contents []entity.ContentForMonitor

	rows, err := r.db.Query(
		"SELECT file_name, file_path, start_time, end_time FROM content WHERE macaddress = $1 AND latest_history = $2",
		MacAddress, 3,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var content entity.ContentForMonitor
		if err := rows.Scan(&content.FileName, &content.FilePath, &content.StartTime, &content.EndTime); err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}


func (r *SocketRepository) SetMonitorStatusActive(MacAddress string) error{
	return nil
}


func (r *SocketRepository) SetMonitorStatus(macAddress string, status bool) error {
	query := `UPDATE monitors SET status = $1 WHERE macaddress = $2`
	_, err := r.db.Exec(query, status, macAddress)
	return err
}
