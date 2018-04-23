package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/wwgberlin/go-weather-widget/weather"
)

type (
	rendererMock struct {
		buildFunc  func(...string) (*template.Template, error)
		renderFunc func(http.ResponseWriter, *template.Template, interface{}) error
		pathFunc   func(s ...string) []string
	}
	forcasterMock struct {
		forecast func(string) (*weather.Conditions, error)
	}
)

func (f forcasterMock) Forecast(s string) (*weather.Conditions, error) {
	return f.forecast(s)
}

func (rdr rendererMock) BuildTemplate(dep ...string) (*template.Template, error) {
	return rdr.buildFunc(dep...)
}

func (rdr rendererMock) RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	return rdr.renderFunc(w, tmpl, data)
}

func (rdr rendererMock) PathToTemplateFiles(s ...string) []string {
	return rdr.pathFunc(s...)
}

func TestIndexHandler_BuildTemplate(t *testing.T) {
	buildInvoked := false

	rdr := rendererMock{
		pathFunc: func(layouts ...string) []string {
			expectedTemplateLayouts := []string{
				"index.tmpl",
				"layouts/layout.tmpl",
				"layouts/head.tmpl",
			}
			if err := checkTemplates(layouts, expectedTemplateLayouts); err != nil {
				t.Error(err)
			}
			return []string{"path/a", "path/b", "path/c"}
		},
		buildFunc: func(layouts ...string) (*template.Template, error) {
			buildInvoked = true

			expectedTemplateLayouts := []string{"path/a", "path/b", "path/c"}

			if err := checkTemplates(layouts, expectedTemplateLayouts); err != nil {
				t.Error(err)
			}
			return template.New("some template"), nil
		},
	}

	indexHandler(rdr)

	if !buildInvoked {
		t.Error("BuildTemplated was expected to be called")
	}
}

func TestIndexHandler_TestRender(t *testing.T) {
	myTmpl := template.New("some template")
	req := httpGetRequest("some_path")
	rr := httptest.NewRecorder()

	rdr := rendererMock{
		pathFunc: func(layouts ...string) []string { return layouts },
		buildFunc: func(layouts ...string) (*template.Template, error) {
			return myTmpl, nil
		},
		renderFunc: func(w http.ResponseWriter, tmpl *template.Template, i interface{}) error {
			if tmpl != myTmpl {
				t.Errorf("Unexpected argument in call to RenderTemplate. Wanted template %s got %s", myTmpl.Name(), tmpl.Name())
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("all good"))
			return nil
		},
	}

	http.HandlerFunc(indexHandler(rdr)).ServeHTTP(rr, req)

	if err := checkResponse(rr.Code, http.StatusOK,
		rr.Body.String(), "all good"); err != nil {
		t.Error(err.Error())
	}
}

func TestIndexHandler_FailToBuildTemplate(t *testing.T) {
	recovered := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
			}
		}()

		indexHandler(rendererMock{
			pathFunc: func(s ...string) []string { return s },
			buildFunc: func(s2 ...string) (*template.Template, error) {
				return nil, errors.New("some error")
			},
		})
	}()
	if !recovered {
		t.Error("Expected to panic on BuildTemplate failure")
	}
}

func TestIndexHandler_FailToRenderTemplate(t *testing.T) {
	req := httpGetRequest("some_path")
	rr := httptest.NewRecorder()
	rdr := rendererMock{
		pathFunc: func(s ...string) []string {
			return s
		},
		buildFunc: func(s2 ...string) (*template.Template, error) {
			return nil, nil
		},
		renderFunc: func(http.ResponseWriter, *template.Template, interface{}) error {
			return errors.New("some error occurred")
		},
	}

	http.HandlerFunc(indexHandler(rdr)).ServeHTTP(rr, req)
	body := rr.Body.String()[0 : len(rr.Body.String())-1]

	if err := checkResponse(
		rr.Code,
		http.StatusInternalServerError,
		body,
		"some error occurred",
	); err != nil {
		t.Error(err.Error())
	}
}

func TestWidgetHandler_TestBuild(t *testing.T) {
	buildInvoked := false

	rdr := rendererMock{
		pathFunc: func(layouts ...string) []string {
			expectedTemplateLayouts := []string{
				"widget.tmpl",
				"layouts/layout.tmpl",
				"layouts/head.tmpl",
			}
			if err := checkTemplates(layouts, expectedTemplateLayouts); err != nil {
				t.Error(err)
			}
			return []string{"path/a", "path/b", "path/c"}
		},
		buildFunc: func(layouts ...string) (*template.Template, error) {
			buildInvoked = true

			expectedTemplateLayouts := []string{"path/a", "path/b", "path/c"}

			if err := checkTemplates(layouts, expectedTemplateLayouts); err != nil {
				t.Error(err)
			}
			return template.New("some template"), nil
		},
	}

	widgetHandler(rdr, forcasterMock{})

	if !buildInvoked {
		t.Error("BuildTemplated was expected to be called")
	}
}

func TestWidgetHandler_TestRender(t *testing.T) {
	myTmpl := template.New("some template")
	req := httpGetRequest("?location=myLocation")
	rr := httptest.NewRecorder()

	conditions := setupConditions()

	forecaster := forcasterMock{
		forecast: func(s string) (*weather.Conditions, error) {
			if s != "myLocation" {
				t.Errorf("Unexpected argument in call to Forecast")
			}
			return &conditions, nil
		},
	}
	rdr := rendererMock{
		pathFunc: func(layouts ...string) []string { return layouts },
		buildFunc: func(layouts ...string) (*template.Template, error) {
			return myTmpl, nil
		},
		renderFunc: func(w http.ResponseWriter, tmpl *template.Template, v interface{}) error {
			if tmpl != myTmpl {
				t.Errorf("Unexpected argument in call to RenderTemplate. Wanted template %s got %s", myTmpl.Name(), tmpl.Name())
			}

			if m, ok := v.(map[string]interface{}); !ok {
				t.Error("Unexpected type in call to RenderTemplate. Want map[string]interface{} but got", reflect.TypeOf(v))
			} else {
				if err := checkMapFields(conditions, m); err != nil {
					t.Error(err)
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("all good"))
			return nil
		},
	}

	http.HandlerFunc(widgetHandler(rdr, forecaster)).ServeHTTP(rr, req)

	if err := checkResponse(rr.Code, http.StatusOK,
		rr.Body.String(), "all good"); err != nil {
		t.Error(err.Error())
	}
}

func TestWidgetHandler_FailToBuildTemplate(t *testing.T) {
	recovered := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = true
			}
		}()

		widgetHandler(rendererMock{
			pathFunc: func(s ...string) []string { return s },
			buildFunc: func(s2 ...string) (*template.Template, error) {
				return nil, errors.New("some error")
			},
		}, forcasterMock{})
	}()
	if !recovered {
		t.Error("Expected to panic on BuildTemplate failure")
	}
}

func httpGetRequest(path string) *http.Request {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func checkTemplates(layouts, expectedLayouts []string) error {
	if !reflect.DeepEqual(expectedLayouts, layouts) {
		return fmt.Errorf(`
				Unexpected arguments in call to BuildTemplate. 
				wanted: %v received %v`,
			expectedLayouts,
			layouts,
		)
	}
	return nil
}

func checkResponse(code, expectedCode int, body, expectedBody string) error {
	if code != expectedCode {
		return fmt.Errorf("handler returned wrong status code: got %d want %d",
			code, expectedCode)
	}

	if body != expectedBody {
		return fmt.Errorf("handler returned unexpected body: got %s want %s",
			body, expectedBody)
	}
	return nil
}

func setupConditions() weather.Conditions {
	return weather.Conditions{
		Location:    "some location",
		Description: "description",
		Celsius:     5,
	}
}

func checkMapFields(conditions weather.Conditions, m map[string]interface{}) error {
	expected := map[string]interface{}{
		"location":    conditions.Location,
		"celsius":     conditions.Celsius,
		"description": conditions.Description,
	}

	if !reflect.DeepEqual(expected, m) {
		return fmt.Errorf("unexpected arguments in call to RenderTemplate. Wanted %v but got %v", expected, m)
	}
	return nil
}
