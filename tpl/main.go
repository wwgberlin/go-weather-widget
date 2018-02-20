package tpl

import (
	"fmt"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, location string, data interface{}) error {
	t, err := template.ParseFiles(fmt.Sprintf("templates/%v.go.html", location))
	if err != nil {
		return nil
	}
	return t.Execute(w, data)
}
