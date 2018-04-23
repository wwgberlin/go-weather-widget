package tpl

import (
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"
)

func TestTemplatesLayout(t *testing.T) {
	expected := "<HTML><BODY></BODY></HTML>"
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Error("layout.tmpl was expected to parse without any errors")
	}
	rr := httptest.NewRecorder()
	tmpl.ExecuteTemplate(rr, "layout", "some data")

	if expected != stripWhitespaces(rr.Body.String()) {
		t.Errorf("Unexpected body rendered from template\nWanted:\n%s\n but got:\n%s", expected, rr.Body.String())
	}

	if tmpl.Lookup("layout") == nil {
		t.Error("layout.tmpl was expected to define template layout")
	}
	if tmpl.Lookup("content") == nil {
		t.Error("layout.tmpl was expected to define empty template content")
	}
	if tmpl.Lookup("head") == nil {
		t.Error("layout.tmpl was expected to define empty template head")
	}
	if tmpl.Lookup("scripts") == nil {
		t.Error("layout.tmpl was expected to define empty template scripts")
	}
	if tmpl.Lookup("styles") == nil {
		t.Error("layout.tmpl was expected to define empty template styles")
	}
}

func TestTemplatesHead(t *testing.T) {
	expected := "<HEAD><TITLE></TITLE></HEAD>"

	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")
	if err != nil {
		t.Error("head.tmpl was expected to parse without any errors")
	}
	rr := httptest.NewRecorder()
	tmpl.ExecuteTemplate(rr, "head", nil)

	if expected != stripWhitespaces(rr.Body.String()) {
		t.Errorf("Unexpected body rendered from template\nWanted:\n%s\n but got:\n%s", expected, rr.Body.String())
	}
}

func TestTemplatesHeadWithTitle(t *testing.T) {
	expected := "<HEAD><TITLE>some_data</TITLE></HEAD>"
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Error("head.tmpl was expected to parse without any errors")
	}

	tmpl, err = tmpl.Parse(`{{define "title"}}{{.}}{{end}}`)
	rr := httptest.NewRecorder()
	tmpl.ExecuteTemplate(rr, "head", "some_data")

	if expected != stripWhitespaces(rr.Body.String()) {
		t.Errorf("Unexpected body rendered from template\nWanted:\n%s\n but got:\n%s", expected, rr.Body.String())
	}
}

func stripWhitespaces(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}
