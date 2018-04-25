package tpl

import (
	"html/template"
	"io"
)

type LayoutRenderer struct {
	helpers    template.FuncMap
	layoutName string
}

func NewRenderer(helpers template.FuncMap, layoutName string) *LayoutRenderer {
	return &LayoutRenderer{
		helpers:    helpers,
		layoutName: layoutName,
	}
}

// BuildTemplate attempts to build a new template given the layoutName,
// the helpers FuncMap defined in the renderer, and parses the files.
// Use template.Must to panic if parse fails.
func (r *LayoutRenderer) BuildTemplate(files ...string) *template.Template {
	return template.Must(template.New(r.layoutName).Funcs(r.helpers).ParseFiles(files...))
}

// RenderTemplate executes the provided template and returns the error
// if the execution fails.
func (r *LayoutRenderer) RenderTemplate(w io.Writer, tmpl *template.Template, data interface{}) error {
	return tmpl.Execute(w, data)
}
