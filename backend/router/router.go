package router

import (
	"database/sql"
	"encoding/json"

	//"fmt"
	"log"
	"net/http"

	"github.com/OmarCodes2/MacShuttle/database"
	"github.com/OmarCodes2/MacShuttle/models"
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
