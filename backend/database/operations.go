package database

import (
	"database/sql"
	"github.com/OmarCodes2/MacShuttle/models"
)

func SaveLocation(db *sql.DB, locData models.LocationData, newRunID int) error {
	query := `INSERT INTO bus_positions (run_id, timestamp_ms, geom, direction) VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)`
	_, err := db.Exec(query, newRunID, locData.Timestamp, locData.Longitude, locData.Latitude, locData.Direction)
	return err
}

func GetNewRunID(db *sql.DB) (int, error) {
	var latestRunID int
	err := db.QueryRow(`SELECT COALESCE(MAX(run_id), 0) FROM bus_positions`).Scan(&latestRunID)
	if err != nil {
		return 0, err
	}
	newRunID := latestRunID + 1

	return newRunID, nil
}