package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ecosia/women-who-go/handlers"
)

func main() {
	port := os.Getenv("PORT")
	http.Handle("/weather", handlers.NewWeatherHandler())
	http.Handle("/", handlers.NewIndexHandler())

	log.Printf("Application serving on http://localhost:%s ...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
