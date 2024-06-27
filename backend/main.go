package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/OmarCodes2/MacShuttle/database"
	"github.com/OmarCodes2/MacShuttle/reference"
	"github.com/OmarCodes2/MacShuttle/router"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func fetchETAData() ([]float64, error) {
	resp, err := http.Get("http://localhost:5000/getETA")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	var data []float64
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Fetch ETAs from the backend route
		data, err := fetchETAData()
		if err != nil {
			log.Println("Error fetching ETA data:", err)
			break
		}

		if len(data) != 4 {
			log.Println("Unexpected data format")
			break
		}

		// Round the first two elements to the nearest whole number
		stop1ETA := math.Ceil(data[0])
		stop2ETA := math.Ceil(data[1])

		// Prepare the data to send through the WebSocket
		dummyData := map[string]interface{}{
			"etas": map[string]string{
				"stop1": fmt.Sprintf("ETA: %.0f minutes", stop1ETA),
				"stop2": fmt.Sprintf("ETA: %.0f minutes", stop2ETA),
			},
			"busPosition": map[string]float64{
				"latitude":  data[3],
				"longitude": data[2],
			},
		}

		err = conn.WriteJSON(dummyData)
		if err != nil {
			log.Println("Error writing JSON:", err)
			break
		}

		// Send updates every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func insertBusPositions(db *sql.DB, runID int) {
	for _, stop := range reference.ReferenceMap {
		log.Printf("running bus position simulator")
		err := database.SaveLocation1(db, stop.Longitude, stop.Latitude, stop.Direction, stop.TimeStamp, runID)
		if err != nil {
			log.Fatalf("Error inserting stop info: %v", err)
		}
		log.Printf("Inserted stop info: %+v", stop)
		time.Sleep(2 * time.Second) // Sleep for 2 seconds
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("POSTGRES_CONNECTION_STRING")

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Successfully connected to the PostgreSQL database")
	//bus poition similator
	log.Println("bus simulation")
	insertBusPositions(db, 23)

	r := router.InitializeRouter(db)
	r.HandleFunc("/ws", wsHandler)
	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
