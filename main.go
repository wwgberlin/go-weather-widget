package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/wwgberlin/go-weather-widget/tpl"
	"github.com/wwgberlin/go-weather-widget/weather/worldweatheronline"
)

func main() {
	const (
		layoutsPath        = "./tpl/templates"
		layoutTemplateName = "layout"
	)

	port := flag.String("port", "8080", "Optional: 4 bytes port")
	apiKey := flag.String("api_key", "", "Required")
	flag.Parse()

	if !validateInput(*port, *apiKey) {
		flag.Usage()
		return
	}

	rdr := tpl.NewRenderer(tpl.Helpers, layoutTemplateName)

	http.HandleFunc("/", indexHandler(layoutsPath, rdr))
	http.HandleFunc("/weather", widgetHandler(layoutsPath, rdr, worldweatheronline.New(*apiKey)))
	http.Handle("/images/", http.StripPrefix("/", http.FileServer(http.Dir("./public/static"))))
	http.Handle("/styles/", http.StripPrefix("/", http.FileServer(http.Dir("./public/static"))))

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
