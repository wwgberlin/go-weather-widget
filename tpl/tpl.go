package tpl

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	path string
}

func NewRenderer(pathToTemplates string) *Renderer {
	return &Renderer{
		path: pathToTemplates,
	}
}

// BuildTemplate takes the name of the template and it's dependencies
// make a slice of type string of size len(dependencies + 1)
// initialize the first item of the slice with the path to the template
// iterate over the dependencies and add the respective file names to the array
func (r *Renderer) BuildTemplate(name string, files ...string) (*template.Template, error) {
	return template.New(name).Funcs(helpers).ParseFiles(files...)
}

func (r *Renderer) RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "layout", data)
}

func (r *Renderer) PathToTemplateFiles(templates ...string) []string {
	files := make([]string, len(templates))

	for i, t := range templates {
		files[i] = filepath.Join(r.path, t)
	}
	return files
}
