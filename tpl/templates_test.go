package tpl

import (
	"bytes"
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
	tmpl.ExecuteTemplate(rr, "layout", "some data")

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
	tmpl.ExecuteTemplate(rr, "layout", "TITLE")
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
	tmpl.ExecuteTemplate(rr, "layout", "content")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))

	if body := doc.Find("body"); body.Length() == 0 {
		t.Error("Expected to render title element")
	} else if strings.TrimSpace(body.Text()) != "content" {
		t.Errorf("head.tmpl was expected to be rendered with body 'content' but got %s", body.Text())
	}
}

func TestTemplateHead(t *testing.T) {
	//expected := "<head></head>"

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
	tmpl.ExecuteTemplate(rr, "head", nil)

	//body := stripWhitespaces(rr.Body.String())
	//if expected != body {
	//	t.Errorf("Unexpected body rendered from template\nWanted:\n%s\n but got:\n%s", expected, body)
	//}
}

func TestTemplatesHeadWithTitle(t *testing.T) {
	//expected := "<head>some_title</head>"
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "title"}}{{.}}{{end}}`)
	rr := httptest.NewRecorder()
	tmpl.ExecuteTemplate(rr, "head", "some_title")

	//body := stripWhitespaces(rr.Body.String())
	//if expected != body {
	//	t.Errorf("Unexpected body rendered from template\nWanted:\n%s\n but got:\n%s", expected, body)
	//}
}

func TestTemplatesHeadWithStyles(t *testing.T) {
	//expected := "<head>some_data</head>"
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "styles"}}{{.}}{{end}}`)

	rr := httptest.NewRecorder()
	tmpl.ExecuteTemplate(rr, "head", "some_data")

	//body := stripWhitespaces(rr.Body.String())
	//if expected != body {
	//	t.Errorf("Unexpected body rendered from template\nWanted:\n%s\n but got:\n%s", expected, body)
	//}
}

func TestTemplateWidget(t *testing.T) {
	tmpl := template.New("widget").Funcs(Helpers)
	tmpl, err := tmpl.ParseFiles("./templates/widget.tmpl")

	if err != nil {
		t.Fatalf("widget.tmpl was expected to parse without any errors. %v", err)
	}

	if tmpl.Lookup("content") == nil {
		t.Error("widget.tmpl was expected to define template content")
	}

	rr := httptest.NewRecorder()

	tmpl.ExecuteTemplate(rr, "content", map[string]interface{}{
		"location":    "Berlin",
		"description": "It's beautiful",
		"celsius":     25,
	})

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rr.Body.Bytes()))
	if doc.Find("div.base").Length() == 0 {
		t.Error("expected to render div with class 'base'")
	} else {
		gopherDiv := doc.Find("div.base")
		if gopherDiv.Find(".sunglasses").Length() == 0 {
			t.Error("gopher was expected to have sunglasses")
		}
	}

}
