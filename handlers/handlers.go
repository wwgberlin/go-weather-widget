package handlers

import (
	"net/http"

	"github.com/wwgberlin/go-weather-widget/tpl"
	"github.com/wwgberlin/go-weather-widget/weather/worldweatheronline"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.RenderTemplate(w, "index.tmpl", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WeatherHandler returns a new handle that is supposed
// to serve the app's result page
func WeatherHandler() func(w http.ResponseWriter, r *http.Request) {
	forecaster := worldweatheronline.New()
	return func(w http.ResponseWriter, r *http.Request) {

		location := r.URL.Query().Get("location")
		if location == "" {
			http.Error(w, "Must specify location", http.StatusBadRequest)
			return
		}

		data, err := forecaster.Forecast(location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tpl.RenderTemplate(w, "weather.tmpl", map[string]interface{}{
			"location":    data.Location(),
			"description": data.Description(),
			"celsius":     data.Celsius(),
			"code":        data.WeatherCode(),
			"query":       location,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
