package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/OmarCodes2/MacShuttle/router"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Dummy ETA data
		dummyData := map[string]string{
			"stop1": "ETA: 2 minutes",
			"stop2": "ETA: 5 minutes",
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

	r := router.InitializeRouter(db)
	r.HandleFunc("/ws", wsHandler)
	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
