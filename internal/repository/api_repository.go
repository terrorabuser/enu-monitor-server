package repository

import (
	"database/sql"
)


type ApiRepository struct {
	db *sql.DB
}

func NewApiRepository(db *sql.DB) *ApiRepository {
	return &ApiRepository{db: db}
}


func (r *ApiRepository) GetBuildings() ([]string, error) {
	var buildings []string

	rows, err := r.db.Query("SELECT DISTINCT building FROM monitors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var building string
		if err := rows.Scan(&building); err != nil {
			return nil, err
		}
		buildings = append(buildings, building)
	}
	return buildings, nil
}


func (r *ApiRepository) GetFloors(building string) ([]string, error) {
	var floors []string

	rows, err := r.db.Query("SELECT DISTINCT floor FROM monitors WHERE building = $1", building)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var floor string
		if err := rows.Scan(&floor); err != nil {
			return nil, err
		}
		floors = append(floors, floor)
	}
	return floors, nil
}


// Получает примечания для конкретного здания и этажа
func (r *ApiRepository) GetNotes(building string, floor int) ([]string, error) {
	rows, err := r.db.Query("SELECT notes FROM monitors WHERE building = $1 AND floor = $2", building, floor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var note string
		if err := rows.Scan(&note); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}


