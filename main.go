package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/wwgberlin/go-weather-widget/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		portPtr := flag.String("port", ":9999", "port")
		flag.Parse()
		port = *portPtr
	}
	http.HandleFunc("/weather", handlers.WeatherHandler())
	http.HandleFunc("/", handlers.IndexHandler)

	log.Printf("Application serving on http://localhost:%s ...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
