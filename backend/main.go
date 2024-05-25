package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OmarCodes2/MacShuttle/router"
)

func main() {
	r := router.InitializeRouter()
	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
