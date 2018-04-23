package tpl

import (
	"bytes"
	"errors"
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestTemplateLayout(t *testing.T) {
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Fatalf("layout.tmpl was expected to parse without any errors. %v", err)
	}

	rr := httptest.NewRecorder()
	if err = tmpl.ExecuteTemplate(rr, "layout", "some data"); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if body := doc.Find("html body"); body.Length() == 0 {
		t.Error("Expected to render html and body elements")
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
}

func TestLayoutWithHead(t *testing.T) {
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "head"}}<head><title>{{.}}</title></head>{{end}}`)

	rr := httptest.NewRecorder()
	if err = tmpl.ExecuteTemplate(rr, "layout", "TITLE"); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if title := doc.Find("title"); title.Length() == 0 {
		t.Error("Expected to render title element")
	} else if title.Text() != "TITLE" {
		t.Errorf("head.tmpl was expected to be rendered with title'TITLE' but got %s", title.Text())
	}
}

func TestLayoutWithContent(t *testing.T) {
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "content"}}{{.}}{{end}}`)

	rr := httptest.NewRecorder()
	if err = tmpl.ExecuteTemplate(rr, "layout", "content"); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if body := doc.Find("body"); body.Length() == 0 {
		t.Error("Expected to render title element")
	} else if strings.TrimSpace(body.Text()) != "content" {
		t.Errorf("head.tmpl was expected to be rendered with body 'content' but got %s", body.Text())
	}
}

func TestTemplateHead(t *testing.T) {
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")
	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	if tmpl.Lookup("styles") == nil {
		t.Error("head.tmpl was expected to define empty template styles")
	}

	if tmpl.Lookup("title") == nil {
		t.Error("head.tmpl was expected to define empty template title")
	}

	rr := httptest.NewRecorder()
	if err = tmpl.ExecuteTemplate(rr, "head", nil); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if head := doc.Find("head"); head.Length() == 0 {
		t.Error("Expected to render title element")
	}
}

func TestTemplatesHeadWithTitle(t *testing.T) {
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "title"}}<title>{{.}}</title>{{end}}`)

	rr := httptest.NewRecorder()
	if tmpl.ExecuteTemplate(rr, "head", "some_title"); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if head := doc.Find("head"); head.Length() == 0 {
		t.Error("Expected to render title element")
	} else if strings.TrimSpace(head.Text()) != "some_title" {
		html, _ := head.Html()
		t.Errorf("head element expected to have title 'some_title', but got '%s'", html)
	}
}

func TestTemplatesHeadWithStyles(t *testing.T) {
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "styles"}}<link rel={{.}}/>{{end}}`)

	rr := httptest.NewRecorder()
	if err = tmpl.ExecuteTemplate(rr, "head", "some_link"); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if head := doc.Find("head"); head.Length() == 0 {
		t.Error("Expected to render title element")
	} else if html, _ := head.Html(); strings.TrimSpace(html) != "<link rel=\"some_link/\"/>" {
		t.Errorf("head element expected to have title 'some_link', but got '%s'", html)
	}
}

func TestTemplateWidget(t *testing.T) {
	h := copyFuncMap(Helpers)
	h["clothe"] = myClothe

	tmpl := template.New("widget").Funcs(h)
	tmpl, err := tmpl.ParseFiles("./templates/widget.tmpl")

	if err != nil {
		t.Fatalf("widget.tmpl was expected to parse without any errors. %v", err)
	}

	if tmpl.Lookup("content") == nil {
		t.Error("widget.tmpl was expected to define template content")
	}

	if tmpl.Lookup("styles") == nil {
		t.Error("widget.tmpl was expected to define template styles")
	}

	if tmpl.Lookup("scripts") == nil {
		t.Error("widget.tmpl was expected to define template scripts")
	}

	rr := httptest.NewRecorder()

	if err = tmpl.ExecuteTemplate(rr, "content", map[string]interface{}{
		"location":    "Berlin",
		"description": "It's spring time",
		"celsius":     25,
	}); err != nil {
		t.Fatalf("template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))
	if doc.Find("div.base").Length() == 0 {
		t.Error("expected to render div with class 'base'")
	} else {
		gopherDiv := doc.Find("div.base")
		if gopherDiv.Find(".sandals").Length() == 0 {
			t.Error("gopher was expected to have sandals")
		}
	}
}

func myClothe(args ...interface{}) ([]string, error) {
	if len(args) < 2 {
		return nil, errors.New("clothe expects 2 arguments to be passed (description, celsius)")
	}
	if desc, ok := args[0].(string); !ok {
		return nil, errors.New("first argument in clothes was expected to be a string (description)")
	} else if desc != "It's spring time" {
		return nil, errors.New("first argument in clothes was expected to be the weather description")
	}
	if celsius, ok := args[1].(int); !ok {
		return nil, errors.New("second argument in clothes was expected to be an integer (celsius)")
	} else if celsius != 25 {
		return nil, errors.New("first argument in clothes was expected to be the weather celsius")
	}
	return []string{"sandals"}, nil
}

func copyFuncMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
