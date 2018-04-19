package tpl

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

// Load templates on program initialisation
func init() {

	//init is called for every import
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "templates/"

	dependencies, err := filepath.Glob(templatesDir + "layouts/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	pages, err := filepath.Glob(templatesDir + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our layouts/ and dependencies/ directories
	for _, page := range pages {
		//nifty trick here
		files := append([]string{page}, dependencies...)
		templates[filepath.Base(page)] = template.Must(
			template.New(page).Funcs(helpers).ParseFiles(files...),
		)
		//template.Must(template.ParseFiles(files...))
	}
}

func RegisterTemplate(name string, dependencies ...string) {
	files := append([]string{name}, dependencies...)
	templates[name] = template.Must(template.New(name).Funcs(helpers).ParseFiles(files...))
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "layout", data)
}

type Renderer struct {
	template     string
	dependencies []string
}

func (r Renderer) RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base.tmpl", data)
}
