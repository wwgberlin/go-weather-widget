package main

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/wwgberlin/go-weather-widget/weather"
)

type (
	renderer interface {
		BuildTemplate(...string) *template.Template
		RenderTemplate(io.Writer, *template.Template, interface{}) error
	}

	forcaster interface {
		Forecast(string) (*weather.Conditions, error)
	}
)

func indexHandler(layoutsPath string, rdr renderer) func(w http.ResponseWriter, r *http.Request) {
	var (
		tmpl *template.Template
	)
	files := pathToTemplateFiles(layoutsPath, "index.tmpl", "layouts/layout.tmpl", "layouts/head.tmpl")

	tmpl = rdr.BuildTemplate(files...)

	return func(w http.ResponseWriter, r *http.Request) {
		queryStr := r.URL.Query().Get("location")

		if err := rdr.RenderTemplate(w, tmpl, map[string]interface{}{
			"location": queryStr,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// WidgetHandler receives a path to the template files
// renderer and a forecaster and returns a function
// with the http.HandleFunc signature.
//
// Similarly to indexHandler above, instantiate your template with the
// template files in the enclosing function (but replace index.tmpl with
// widget.tmpl).
//
// The forecaster provides a function Forecast that receives a location (string)
// and returns weather.Conditions object with the fields:
// Description, Location, Celsius
// Call forcaster.Forecast with your request param -
// r.URL.Query().Get("location")
//
// Instantiate a map[string]interface{} to pass to template execution
// and add the data from the Conditions to the map (use lower case)
// e.g. m["location"] = c.Location, etc.

func widgetHandler(layoutsPath string, rdr renderer, forecaster forcaster) func(w http.ResponseWriter, r *http.Request) {
	var (
		tmpl *template.Template
	)
	files := pathToTemplateFiles(layoutsPath, "widget.tmpl", "layouts/layout.tmpl", "layouts/head.tmpl")
	tmpl = rdr.BuildTemplate(files...)

	return func(w http.ResponseWriter, r *http.Request) {
		queryStr := r.URL.Query().Get("location")
		data, err := forecaster.Forecast(queryStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := rdr.RenderTemplate(w, tmpl, map[string]interface{}{
			"location":    data.Location,
			"celsius":     data.Celsius,
			"description": data.Description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func pathToTemplateFiles(path string, templates ...string) []string {
	files := make([]string, len(templates))

	for i, t := range templates {
		files[i] = filepath.Join(path, t)
	}
	return files
}
