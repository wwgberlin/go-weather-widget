package handlers

import (
	"net/http"

	"github.com/ecosia/women-who-go/tpl"
	"github.com/ecosia/women-who-go/weather"
	"github.com/ecosia/women-who-go/weather/worldweatheronline"
)

// NewWeatherHandler returns a new handle that is supposed
// to serve the app's result page
func NewWeatherHandler() http.Handler {
	return &weatherHandler{forecaster: worldweatheronline.New()}
}

type weatherHandler struct {
	forecaster weather.Forecaster
}

func (h *weatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if q := r.URL.Query().Get("location"); q != "" {
		data, err := h.forecaster.Forecast(r.URL.Query().Get("location"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Render(w, "weather", map[string]interface{}{
			"location":    data.Location(),
			"description": data.Description(),
			"celsius":     data.Celsius(),
			"query":       q,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.NotFound(w, r)
}
