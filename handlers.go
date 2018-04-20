package main

import (
	"fmt"
	"net/http"

	"github.com/wwgberlin/go-weather-widget/tpl"
	"github.com/wwgberlin/go-weather-widget/weather/worldweatheronline"
)

func whoops(rdr tpl.Renderer, errorCode int, err error, w http.ResponseWriter) {

	if err := rdr.RenderTemplate(w, "whoops.tmpl", map[string]interface{}{"Error": err}); err != nil {
		w.Write([]byte("Something went wrong"))
	}
}

func indexHandler(rdr tpl.Renderer) func(w http.ResponseWriter, r *http.Request) {
	if err := rdr.RegisterTemplate("index.tmpl", "layouts/layout.tmpl", "layouts/head.tmpl"); err != nil {
		panic(err)
	}

	//handle the request
	return func(w http.ResponseWriter, r *http.Request) {
		if err := rdr.RenderTemplate(w, "index.tmpl", nil); err != nil {
			whoops(rdr, http.StatusInternalServerError, err, w)
		}
	}
}

func weatherHandler(rdr tpl.Renderer, apiKey string) func(w http.ResponseWriter, r *http.Request) {
	forecaster := worldweatheronline.New(apiKey)
	if err := rdr.RegisterTemplate("widget.tmpl", "layouts/layout.tmpl", "layouts/head.tmpl"); err != nil {
		panic(err)
	}

	//handle the request
	return func(w http.ResponseWriter, r *http.Request) {
		location := r.URL.Query().Get("location")
		if location == "" {
			whoops(rdr, http.StatusBadRequest, fmt.Errorf("must specify location"), w)
			return
		}

		data, err := forecaster.Forecast(location)
		if err != nil {
			whoops(rdr, http.StatusInternalServerError, err, w)
			return
		}

		if err := rdr.RenderTemplate(w, "widget.tmpl", map[string]interface{}{
			"location":    data.Location(),
			"description": data.Description(),
			"celsius":     data.Celsius(),
			"query":       location,
		}); err != nil {
			whoops(rdr, http.StatusInternalServerError, err, w)
		}
		return
	}
}
