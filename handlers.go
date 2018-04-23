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
		queryStr := r.URL.Query().Get("query")

		if err := rdr.RenderTemplate(w, tmpl, map[string]interface{}{"query": queryStr}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func widgetHandler(rdr renderer, forecaster forcaster) func(w http.ResponseWriter, r *http.Request) {
	var (
		tmpl *template.Template
		err  error
	)

	files := []string{"widget.tmpl", "layouts/layout.tmpl", "layouts/head.tmpl"}
	if tmpl, err = rdr.BuildTemplate(rdr.PathToTemplateFiles(files...)...); err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := forecaster.Forecast(r.URL.Query().Get("location"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := rdr.RenderTemplate(w, tmpl, map[string]interface{}{
			"location":    data.Location,
			"description": data.Description,
			"celsius":     data.Celsius,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
