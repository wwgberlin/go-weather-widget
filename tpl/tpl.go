package tpl

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	path      string
	templates map[string]*template.Template
}

func New(pathToTemplates string) Renderer {
	return Renderer{
		path:      pathToTemplates,
		templates: make(map[string]*template.Template),
	}
}

func (r Renderer) pathToTemplate(template string) string {
	return filepath.Join(r.path, template)
}

// RegisterTemplate takes the name of the template and it's dependencies
func (r Renderer) RegisterTemplate(name string, dependencies ...string) (err error) {
	files := make([]string, len(dependencies)+1)
	files[0] = r.pathToTemplate(name)
	for i, d := range dependencies {
		files[i+1] = r.pathToTemplate(d)
	}

	//todo: should this be here on happen on each request?
	r.templates[name], err = template.New(name).Funcs(helpers).ParseFiles(files...)
	return
}

func (r Renderer) RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	tmpl, ok := r.templates[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "layout", data)
}
