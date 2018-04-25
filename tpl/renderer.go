package tpl

import (
	"errors"
	"html/template"
	"io"
	"strings"
)

var DefaultHelpers = template.FuncMap{
	"title": strings.Title,
}

type LayoutRenderer struct {
	Helpers    template.FuncMap
	LayoutName string
}

func NewRenderer(layoutName string) *LayoutRenderer {
	return &LayoutRenderer{
		Helpers:    DefaultHelpers,
		LayoutName: layoutName,
	}
}

// BuildTemplate attempts to build a new template given the LayoutName,
// the Helpers FuncMap defined in the renderer, and parses the files.
// Use template.Must to panic if parse fails.
func (r *LayoutRenderer) BuildTemplate(files ...string) *template.Template {
	return nil
}

// RenderTemplate executes the provided template and returns the error
// if the execution fails.
func (r *LayoutRenderer) RenderTemplate(w io.Writer, tmpl *template.Template, data interface{}) error {
	return errors.New("not implemented")
}
