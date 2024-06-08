package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/OmarCodes2/MacShuttle/router"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	fmt.Println("Server is running on port 5000")

	req, err := http.NewRequest("GET", "/getETA", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        log.Fatal("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

	fmt.Println("Response of eta handler:", rr.Body.String())

	log.Fatal(http.ListenAndServe(":5000", r))	
}
