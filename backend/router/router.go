package router

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math"

	//"fmt"
	"log"
	"net/http"

	"github.com/OmarCodes2/MacShuttle/database"
	"github.com/OmarCodes2/MacShuttle/models"
	"github.com/OmarCodes2/MacShuttle/reference"
	_ "github.com/lib/pq"
)

var (
	runID int
)

func InitializeRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/startTracking", runIDHandler(db))
	mux.HandleFunc("/liveTracking", locationHandler(db))
	mux.HandleFunc("/getETA", getETAHandler(db))
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func runIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		newRunID, err := database.GetNewRunID(db)
		if err != nil {
			http.Error(w, "Error retrieving current run_id", http.StatusInternalServerError)
			return
		}
		//Modifying Run ID Variable
		runID = newRunID
		log.Println("getting new run id correctly")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"run_id": newRunID})
	}
}

func locationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var locData models.LocationData
		err := json.NewDecoder(r.Body).Decode(&locData)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		log.Printf("Received location: %+v\n", locData)
		log.Printf("RunID is %v", runID)
		log.Printf("Direction is %v", locData.Direction)

		if err := database.SaveLocation(db, locData, runID); err != nil {
			log.Printf("Error saving location: %v\n", err)
			http.Error(w, "Error saving location", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Location received and database connection is successful"))
	}
}

// Using Haversine formula to calculate distance between 2 coords
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var (
		r    = 6371 // Earth radius in kilometers
		dLat = (lat2 - lat1) * (math.Pi / 180.0)
		dLon = (lon2 - lon1) * (math.Pi / 180.0)
		a    = math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*math.Sin(dLon/2)*math.Sin(dLon/2)
		c    = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	)
	distance := float64(r) * c
	return distance // Distance in kilometers
}

func GetClosestStop(lat, lon float64, direction string) (reference.StopInfo, error) {
	var closestStop reference.StopInfo
	var minDistance float64

	minDistance = math.MaxFloat64 // Assigning Min Distance Inf Value

	// Loop through the points and identify the closest stop from reference table to the bus coords
	for _, stop := range reference.ReferenceMap {
		distance := Haversine(lat, lon, stop.Latitude, stop.Longitude)
		// If new minimum distance is found in the same direction, update the minimum reference coord to this new coord
		if distance < minDistance && direction == stop.Direction {
			minDistance = distance
			closestStop = stop
		}
	}

	if minDistance == math.MaxFloat64 {
		return closestStop, errors.New("no stop found in GetClosestStop()")
	}
	return closestStop, nil
}

func CalculateETA(referenceCoords reference.StopInfo) ([]float64, error) {
	const (
		millisInMinute = 60000 // 60,000 milliseconds in a minute for conversion
	)
	var ETAStopA int
	var ETAStopB int
	// Calculating eta when bus is driving from A -> B
	if referenceCoords.Direction == "forward" {
		ETAStopB = reference.StopBtime - referenceCoords.TimeStamp
		ETAStopA = ETAStopB + (reference.StopAtime - reference.StopBtime)
	} else { // Calculating eta when bus is driving from B -> A
		ETAStopA = reference.StopAtime - referenceCoords.TimeStamp
		ETAStopB = ETAStopA + reference.StopBtime
	}

	// Converting to ETA in minutes
	ETAStopA = ETAStopA / millisInMinute
	ETAStopB = ETAStopB / millisInMinute

	if ETAStopA < 0 || ETAStopB < 0 {
		return nil, errors.New("failed to calculate ETA in CalculateETA()")
	}

	ETAs := []float64{float64(ETAStopA), float64(ETAStopB), referenceCoords.Longitude, referenceCoords.Latitude}
	return ETAs, nil
}

func GetBusETA(db *sql.DB) ([]float64, error) {
	//Retrieve most recent bus location from db
	busLocation, err := database.GetLatestBusLocation(db)
	if err != nil {
		return nil, err
	}

	//Get closest reference stop to bus ETA
	closestStop, err := GetClosestStop(busLocation.Latitude, busLocation.Longitude, busLocation.Direction)
	if err != nil {
		return nil, err
	}

	//Calculate ETA to each of the stops
	ETAs, err := CalculateETA(closestStop)
	if err != nil {
		return nil, err
	}

	return ETAs, nil
}

func GetETAHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ETAs, err := GetBusETA(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

		log.Println("ETAs array corresponds to ETAStopA, ETAStopB, BusLongitude, BusLatitude: ")
		log.Println(ETAs)
        response := ETAs
        json.NewEncoder(w).Encode(response)
    }
}
	