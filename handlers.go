package main

import (
	"html/template"
	"net/http"

	"github.com/wwgberlin/go-weather-widget/weather"
)

type (
	renderer interface {
		BuildTemplate(...string) (*template.Template, error)
		RenderTemplate(http.ResponseWriter, *template.Template, interface{}) error
		PathToTemplateFiles(templates ...string) []string
	}

	forcaster interface {
		Forecast(string) (*weather.Conditions, error)
	}
)

func indexHandler(rdr renderer) func(w http.ResponseWriter, r *http.Request) {
	var (
		tmpl *template.Template
		err  error
	)
	files := []string{"index.tmpl", "layouts/layout.tmpl", "layouts/head.tmpl"}
	if tmpl, err = rdr.BuildTemplate(rdr.PathToTemplateFiles(files...)...); err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		queryStr := r.URL.Query().Get("location")

		if err := rdr.RenderTemplate(w, tmpl, map[string]interface{}{"location": queryStr}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WidgetHandler receives a renderer and a forecaster and returns a function
// with the http.HandleFunc signature.
//
// Similarly to indexHandler above, instantiate your template in the enclosing function,
// but replace index.tmpl with widget.tmpl use panic to prevent the app from
// starting with invalid template files.
//
// The forecaster has a function Forecast that receives a location (string)
// and returns weather.Conditions object with the fields:
// Description, Location, Celsius
//
// Call Forecast with your request param - r.URL.Query().Get("location")
// Instantiate a map[string]interface{} and add the data from the Conditions
// to the map (use lower case) e.g. m["location"] = c.Location

func widgetHandler(rdr renderer, forecaster forcaster) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not implemented", http.StatusNotImplemented)
	}
}
