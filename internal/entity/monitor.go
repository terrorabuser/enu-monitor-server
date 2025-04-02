//Модель монитора (ID, Name, GroupID, Status, IP)
package entity

import (
	"database/sql"
)

type Monitor struct {
	MacAddress              string `json:"macaddress"`
	Password			   string `json:"password"`
	PasswordHash		string `json:"password_hash"`
	Building          string `json:"building"`
	Floor             string `json:"floor"`
	Note             string `json:"note"`
	MonitorResolution string `json:"monitor_resolution"`
	Status 		  bool `json:"status"`
	IP                sql.NullString `json:"ip"`
	Power 		   bool `json:"power"`
	LastLog 		 *string `json:"last_log"`
}






