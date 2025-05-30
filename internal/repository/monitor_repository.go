package repository


import ( 
	"golang_gpt/internal/entity"
	"database/sql"
)


type MonitorRepository struct {
	db *sql.DB
}

func NewMonitorRepository(db *sql.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

// Получение мак адреса по зданию и этажу
func (r *ContentRepository) GetMacAddressByLocation(building string, floor int, notes string) (string, error) {
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


func (r *MonitorRepository) GetAllMonitors() ([]entity.Monitor, error){
	rows, err := r.db.Query("SELECT * FROM monitors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []entity.Monitor
    for rows.Next() {
        var m entity.Monitor
        err := rows.Scan(&m.MacAddress, &m.Building, &m.Floor, &m.MonitorResolution, &m.Status, &m.IP, &m.Power, &m.LastLog)
        if err != nil {
            return nil, err
        }
        monitors = append(monitors, m)
    }
    return monitors, nil
}


func (r *MonitorRepository) CheckMonitorByPassword(macaddress string) (*entity.Monitor, error) {
    return &entity.Monitor{}, nil // Возвращаем пустую структуру
}
