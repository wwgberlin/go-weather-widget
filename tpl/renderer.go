package tpl

import (
	"errors"
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

// BuildTemplate attempts to build a new template given the layoutName,
// the helpers FuncMap defined in the renderer, and parses the files
func (r *Renderer) BuildTemplate(files ...string) (*template.Template, error) {
	return nil, errors.New("not implemented")
}

// RenderTemplate executes the provided template by the renderer's layout name
// and returns the error if the execution fails.
func (r *Renderer) RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	return errors.New("not implemented")
}

func (r *Renderer) PathToTemplateFiles(templates ...string) []string {
	files := make([]string, len(templates))

	for i, t := range templates {
		files[i] = filepath.Join(r.path, t)
	}
	return files
}
