package tpl_test

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"testing"

	"github.com/wwgberlin/go-weather-widget/tpl"
)

func renderer(layoutName string, helpers *template.FuncMap) (rdr *tpl.LayoutRenderer) {
	if helpers == nil {
		helpers = &tpl.Helpers
	}
	return tpl.NewRenderer(*helpers, layoutName)
}

func TestBuildTemplate(t *testing.T) {
	rdr := renderer("layout", nil)

	tmpl := rdr.BuildTemplate("./test/success1.tmpl")

	if tmpl == nil {
		t.Error("BuildTemplate returned nil template")
	} else {
		if tmpl.Lookup("success1") == nil {
			t.Error("Template success1 was not found in template. Did you call parse files?")
		}
	}
}

func TestBuildTemplate_FuncMap(t *testing.T) {
	rdr := renderer("doesn't matter", &template.FuncMap{
		"defined": func(v interface{}) interface{} {
			return v
		},
	})

	files := []string{"./test/success1.tmpl", "./test/success2.tmpl"}
	recovered := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
			}
			rdr.BuildTemplate(files...)
		}()
	}()
}

func TestBuildTemplate_Errors(t *testing.T) {
	rdr := renderer("doesn't matter", nil)

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

	rdr := renderer("success1", &template.FuncMap{
		"defined": func(v interface{}) interface{} {
			return v
		},
	})

	files := []string{"./test/success1.tmpl", "./test/success2.tmpl"}
	tmpl := rdr.BuildTemplate(files...)
	rdr.RenderTemplate(&b, tmpl, expected)

	res := strings.TrimSpace(b.String())

	if res != expected {
		t.Errorf("Expected to render '%s' but received '%s'", expected, res)
	}
}

func TestRenderTemplate_ErrorHandling(t *testing.T) {
	var b bytes.Buffer

	rdr := renderer("success1", &template.FuncMap{
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
