package tpl_test

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"testing"

	. "github.com/wwgberlin/go-weather-widget/tpl"
)

func testRenderer(layoutName string, helpers template.FuncMap) (rdr *LayoutRenderer) {
	if helpers == nil {
		helpers = DefaultHelpers
	}
	return &LayoutRenderer{
		Helpers:    helpers,
		LayoutName: layoutName,
	}
}

func TestBuildTemplate(t *testing.T) {
	rdr := testRenderer("layout", nil)

	tmpl := rdr.BuildTemplate("./test/success1.tmpl")

	if tmpl == nil {
		t.Error("BuildTemplate returned nil template")
	} else {
		if tmpl.Lookup("success1") == nil {
			t.Error("Template success1 was not found in template. Did you call ParseFiles?")
		}
	}
}

func TestBuildTemplate_FuncMap(t *testing.T) {
	rdr := testRenderer("doesn't matter", template.FuncMap{
		"defined": func(v interface{}) interface{} {
			return v
		},
	})

	files := []string{"./test/success1.tmpl", "./test/success2.tmpl"}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatal("BuildTemplate was not expected to panic. Did you call Funcs()?")
			}
		}()
		rdr.BuildTemplate(files...)
	}()
}

func TestBuildTemplate_Errors(t *testing.T) {
	rdr := testRenderer("doesn't matter", nil)

	recovered := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
			}
		}()
		rdr.BuildTemplate("./test/fail.tmpl")
	}()

	if !recovered {
		t.Error("BuildTemplate was expected to panic")
	}
}

func TestRenderTemplate(t *testing.T) {
	var b bytes.Buffer
	const expected = "SOME TEXT"

	rdr := testRenderer("success1", template.FuncMap{
		"defined": func(v interface{}) interface{} {
			return v
		},
	})

	files := []string{"./test/success1.tmpl", "./test/success2.tmpl"}
	tmpl := rdr.BuildTemplate(files...)

	if err := rdr.RenderTemplate(&b, tmpl, expected); err != nil {
		t.Fatalf("RenderTemplate was expected to succeed without errors. %v", err)
	}

	res := strings.TrimSpace(b.String())
	if res != expected {
		t.Errorf("Expected to render '%s' but received '%s'", expected, res)
	}
}

func TestRenderTemplate_ErrorHandling(t *testing.T) {
	var b bytes.Buffer

	rdr := testRenderer("success1", template.FuncMap{
		"defined": func(v interface{}) (interface{}, error) {
			return nil, errors.New("some error")
		},
	})

	files := []string{"./test/success1.tmpl", "./test/success2.tmpl"}

	tmpl := rdr.BuildTemplate(files...)

	if err := rdr.RenderTemplate(&b, tmpl, 1); err == nil {
		t.Error("RenderTemplate was expected to return an error")
	}
}
