package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/wwgberlin/go-weather-widget/tpl"
	"github.com/wwgberlin/go-weather-widget/weather/worldweatheronline"
)

func main() {
	port := flag.String("port", "8080", "Optional: 4 bytes port")
	apiKey := flag.String("api_key", "", "Required")
	flag.Parse()

	if !validateInput(*port, *apiKey) {
		flag.Usage()
		return
	}

	rdr := tpl.NewRenderer(filepath.Join(getPath(), "templates"))

	http.HandleFunc("/", indexHandler(rdr))
	http.HandleFunc("/weather", widgetHandler(rdr, worldweatheronline.New(*apiKey)))
	http.Handle("/images/", http.StripPrefix("/", http.FileServer(http.Dir("./public/static"))))

	log.Printf("Application serving on http://localhost:%s ...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

func validateInput(port string, apiKey string) bool {
	var p int16
	if _, err := fmt.Sscanf(port, "%d", &p); err != nil {
		return false
	}

	if apiKey == "" {
		return false
	}

	return true
}

func getPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}
