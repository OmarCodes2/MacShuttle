package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OmarCodes2/MacShuttle/router"
)

func main() {
	r := router.InitializeRouter()
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
