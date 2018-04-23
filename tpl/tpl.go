package tpl

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	path       string
	helpers    template.FuncMap
	layoutName string
}

func NewRenderer(pathToTemplates string, helpers template.FuncMap, layoutName string) *Renderer {
	return &Renderer{
		path:       pathToTemplates,
		helpers:    helpers,
		layoutName: layoutName,
	}
}

// BuildTemplate attempts to build a new template with the given name,
// given files and the helper functions in the renderer.
// it returns the template or the error.
func (r *Renderer) BuildTemplate(files ...string) (*template.Template, error) {
	//return nil, errors.New("not implemented")
	return template.New(r.layoutName).Funcs(r.helpers).ParseFiles(files...)
}

// RenderTemplate executes the provided template by the renderer's layout name
// and returns the error if the execution fails.
func (r *Renderer) RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	//return errors.New("not implemented")
	return tmpl.ExecuteTemplate(w, r.layoutName, data)
}

func (r *Renderer) PathToTemplateFiles(templates ...string) []string {
	files := make([]string, len(templates))

	for i, t := range templates {
		files[i] = filepath.Join(r.path, t)
	}
	return files
}
