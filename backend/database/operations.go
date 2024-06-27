package database

import (
	"database/sql"
	"log"

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

func GetLatestBusLocation(db *sql.DB) (models.LocationData, error) {
    var location models.LocationData
	query := `
	SELECT timestamp_ms, ST_X(geom) as longitude, ST_Y(geom) as latitude, direction 
	FROM bus_positions 
	WHERE run_id = (
		SELECT COALESCE(MAX(run_id), 0) FROM bus_positions
	)
	ORDER BY timestamp_ms DESC 
	LIMIT 1`
	
    err := db.QueryRow(query).Scan(&location.Timestamp, &location.Longitude, &location.Latitude, &location.Direction)
    if err != nil {
        return location, err
    }
	log.Println("DB query returned: ", location)
    return location, nil
}

func SaveLocation1(db *sql.DB,longitude float64, latitude float64, direction string, timestamp int , newRunID int) error {
	query := `INSERT INTO bus_positions (run_id, timestamp_ms, geom, direction) VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)`
	_, err := db.Exec(query, newRunID, timestamp, longitude, latitude, direction)
	return err
}