package router

import (
	"encoding/json"
	"log"
	"net/http"
)

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
	Direction string  `json:"direction"`
}

func InitializeRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/location", locationHandler)
	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var locData LocationData
	err := json.NewDecoder(r.Body).Decode(&locData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received location: %+v\n", locData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Location received"))
}
