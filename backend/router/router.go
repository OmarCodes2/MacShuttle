package router

import (
	"database/sql"
	"encoding/json"

	//"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/OmarCodes2/MacShuttle/models"
	"github.com/OmarCodes2/MacShuttle/database"
)

func InitializeRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/location", locationHandler(db))
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
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

		if err := database.SaveLocation(db, locData, 1); err != nil {
			log.Printf("Error saving location: %v\n", err)
			http.Error(w, "Error saving location", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Location received and database connection is successful"))
	}
}

// func startRunHandler(w http.ResponseWriter, r *http.Request) {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	var currentRunID int
// 	err := db.QueryRow("SELECT current_run_id FROM run_info LIMIT 1").Scan(&currentRunID)
// 	if err != nil {
// 		http.Error(w, "Error retrieving current run_id", http.StatusInternalServerError)
// 		return
// 	}

// 	newRunID := currentRunID + 1

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(fmt.Sprintf(`{"run_id": %d}`, newRunID)))
// }
