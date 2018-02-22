package handlers

import (
	"net/http"

	"github.com/ecosia/women-who-go/tpl"
	"github.com/ecosia/women-who-go/weather"
	"github.com/ecosia/women-who-go/weather/worldweatheronline"
)

// NewIndexHandler returns a new handle that is supposed
// to serve the app's index page
func NewIndexHandler() http.Handler {
	return &indexHandler{forecaster: worldweatheronline.New()}
}

type indexHandler struct {
	forecaster weather.Forecaster
}

func (i *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if q := r.URL.Query().Get("location"); q != "" {
		data, err := i.forecaster.Forecast(r.URL.Query().Get("location"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tpl.Render(w, "index", map[string]interface{}{
			"location":    data.Location(),
			"description": data.Description(),
			"celsius":     data.Celsius(),
			"query":       q,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	tpl.Render(w, "index", map[string]interface{}{"q": nil})
}
