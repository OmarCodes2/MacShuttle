package router

import (
	"encoding/json"
	"log"
	"net/http"
)

type LocationData struct {
	Type       string  `json:"_type"`
	SSID       string  `json:"SSID"`
	Altitude   float64 `json:"alt"`
	Battery    int     `json:"batt"`
	Connection string  `json:"conn"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
	Timestamp  int64   `json:"tst"`
	Velocity   float64 `json:"vel"`
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
