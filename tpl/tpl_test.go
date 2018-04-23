package tpl_test

import (
	"errors"
	"html/template"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/wwgberlin/go-weather-widget/tpl"
)

func renderer(layoutName string, helpers *template.FuncMap) (rdr *tpl.Renderer) {
	if helpers == nil {
		helpers = &tpl.Helpers
	}
	return tpl.NewRenderer("./test", *helpers, layoutName)
}

func TestBuildTemplate(t *testing.T) {
	rdr := renderer("layout", nil)

	tmplName := "myName"
	files := rdr.PathToTemplateFiles("success1.tmpl")
	tmpl, err := rdr.BuildTemplate(tmplName, files...)

	if err != nil {
		t.Errorf("Unexpected error received %v", err)
	}
	if tmpl.Name() != "myName" {
		t.Errorf("Unexpected template name. Wanted %s but got %s", tmplName, tmpl.Name())
	}
	if tmpl.Lookup("success1") == nil {
		t.Error("Template success1 was not found in template. Did you call parse files?")
	}
}

func TestBuildTemplate_FuncMap(t *testing.T) {
	rdr := renderer("", &template.FuncMap{
		"defined": func(v interface{}) interface{} {
			return v
		},
	})

	files := rdr.PathToTemplateFiles("success1.tmpl", "success2.tmpl")
	if _, err := rdr.BuildTemplate("some name", files...); err != nil {
		t.Errorf("Unexpected error received %v. Did you call FuncMap?", err)
	}
}

func TestBuildTemplate_Errors(t *testing.T) {
	rdr := renderer("", nil)

	files := rdr.PathToTemplateFiles("fail.tmpl")
	if _, err := rdr.BuildTemplate("some name", files...); err == nil {
		t.Error("BuildTemplate was expected to return an error")
	}
}

var expected = regexp.MustCompile("SOME TEXT\\s+1")

func TestRenderTemplate(t *testing.T) {
	rdr := renderer("success1", &template.FuncMap{
		"defined": func(v interface{}) interface{} {
			return v
		},
	})

	files := rdr.PathToTemplateFiles("success1.tmpl", "success2.tmpl")

	tmpl, err := rdr.BuildTemplate("some templates", files...)
	if err != nil {
		t.Errorf("Unexpected error received %v", err)
	}

	rr := httptest.NewRecorder()
	rdr.RenderTemplate(rr, tmpl, 1)

	if !expected.Match(rr.Body.Bytes()) {
		t.Errorf("%s expected ", rr.Body.String())
	}
}

func TestRenderTemplate_ErrorHandling(t *testing.T) {
	rdr := renderer("success1", &template.FuncMap{
		"defined": func(v interface{}) (interface{}, error) {
			return nil, errors.New("some error")
		},
	})

	files := rdr.PathToTemplateFiles("success1.tmpl", "success2.tmpl")

	tmpl, err := rdr.BuildTemplate("some templates", files...)
	if err != nil {
		t.Errorf("Unexpected error received when calling BuildTemplate %v", err)
	}

	rr := httptest.NewRecorder()
	err = rdr.RenderTemplate(rr, tmpl, 1)

	if err == nil {
		t.Error("RenderTemplate was expected to return an error")
	}
}

func TestPathToTemplateFiles(t *testing.T) {
	rdr := tpl.NewRenderer("some_path", nil, "")
	expected := []string{"some_path/a", "some_path/b"}
	returned := rdr.PathToTemplateFiles("a", "b")
	if !reflect.DeepEqual(expected, returned) {
		t.Errorf("PathToTemplateFiles returned incorrect result. Wanted %v but got %v", expected, returned)
	}
}
