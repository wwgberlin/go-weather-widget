package tpl

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
)

// Render looks up the given location and renders the template against
// the given data, writing to the passed io.Writer
func Render(w io.Writer, location string, data interface{}) error {
	b, err := ioutil.ReadFile(fmt.Sprintf("templates/%s.go.html", location))
	if err != nil {
		return err
	}
	t, err := template.New(location).Funcs(helpers).Parse(string(b))
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}
