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
		buildInvoked  bool
		renderInvoked bool
		pathInvoked   bool
		buildFunc     func(...string) (*template.Template, error)
		renderFunc    func(http.ResponseWriter, *template.Template, interface{}) error
		pathFunc      func(s ...string) []string
	}
	forcasterMock struct {
		invoked  bool
		forecast func(string) (*weather.Conditions, error)
	}
)

func (f forcasterMock) Forecast(s string) (*weather.Conditions, error) {
	return f.forecast(s)
}

func (rdr *rendererMock) BuildTemplate(dep ...string) (*template.Template, error) {
	rdr.buildInvoked = true
	return rdr.buildFunc(dep...)
}

func (rdr *rendererMock) RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) error {
	rdr.renderInvoked = true
	return rdr.renderFunc(w, tmpl, data)
}

func (rdr *rendererMock) PathToTemplateFiles(s ...string) []string {
	rdr.pathInvoked = true
	return rdr.pathFunc(s...)
}

func TestIndexHandler_BuildTemplate(t *testing.T) {
	expectedInputForPath := []string{
		"index.tmpl",
		"layouts/layout.tmpl",
		"layouts/head.tmpl",
	}

	expectedInputForBuild := []string{"path/a", "path/b", "path/c"}
	rdr := &rendererMock{
		pathFunc: func(layouts ...string) []string {
			if err := checkTemplates(layouts, expectedInputForPath); err != nil {
				t.Error(err)
			}
			return expectedInputForBuild
		},
		buildFunc: func(layouts ...string) (*template.Template, error) {
			if err := checkTemplates(layouts, expectedInputForBuild); err != nil {
				t.Error(err)
			}
			return template.New("some template"), nil
		},
	}

	indexHandler(rdr)

	if !rdr.buildInvoked {
		t.Error("BuildTemplated was expected to be called")
	}
	if !rdr.pathInvoked {
		t.Error("PathToTemplateFiles was expected to be called")
	}
}

func TestIndexHandler_TestRender(t *testing.T) {
	const (
		expectedResult = "all good"
	)

	myTmpl := template.New("some template")
	req := httpGetRequest("some_path")
	rr := httptest.NewRecorder()

	rdr := &rendererMock{
		pathFunc: func(layouts ...string) []string { return layouts },
		buildFunc: func(layouts ...string) (*template.Template, error) {
			return myTmpl, nil
		},
		renderFunc: func(w http.ResponseWriter, tmpl *template.Template, i interface{}) error {
			if tmpl != myTmpl {
				t.Errorf("Unexpected argument in call to RenderTemplate. Wanted template %s got %s", myTmpl.Name(), tmpl.Name())
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedResult))
			return nil
		},
	}

	http.HandlerFunc(indexHandler(rdr)).ServeHTTP(rr, req)

	if err := checkResponse(rr.Code, http.StatusOK,
		rr.Body.String(), expectedResult); err != nil {
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

		indexHandler(&rendererMock{
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
	const errMsg = "some error occurred"

	req := httpGetRequest("some_path")
	rr := httptest.NewRecorder()
	rdr := &rendererMock{
		pathFunc: func(s ...string) []string {
			return s
		},
		buildFunc: func(s2 ...string) (*template.Template, error) {
			return nil, nil
		},
		renderFunc: func(http.ResponseWriter, *template.Template, interface{}) error {
			return errors.New(errMsg)
		},
	}

	http.HandlerFunc(indexHandler(rdr)).ServeHTTP(rr, req)
	body := rr.Body.String()[0 : len(rr.Body.String())-1]

	if err := checkResponse(
		rr.Code,
		http.StatusInternalServerError,
		body,
		errMsg,
	); err != nil {
		t.Error(err.Error())
	}
}

func TestWidgetHandler_TestBuild(t *testing.T) {
	expectedInputForPath := []string{
		"widget.tmpl",
		"layouts/layout.tmpl",
		"layouts/head.tmpl",
	}

	expectedInputForBuild := []string{"path/a", "path/b", "path/c"}

	rdr := &rendererMock{
		pathFunc: func(layouts ...string) []string {
			if err := checkTemplates(layouts, expectedInputForPath); err != nil {
				t.Error(err)
			}
			return expectedInputForBuild
		},
		buildFunc: func(layouts ...string) (*template.Template, error) {
			if err := checkTemplates(layouts, expectedInputForBuild); err != nil {
				t.Error(err)
			}
			return template.New("some template"), nil
		},
	}

	widgetHandler(rdr, forcasterMock{})

	if !rdr.pathInvoked {
		t.Error("PathToTemplateFiles was expected to be called")
	}
	if !rdr.buildInvoked {
		t.Error("BuildTemplated was expected to be called")
	}
}

func TestWidgetHandler_TestRender(t *testing.T) {
	const (
		queryLocation  = "myLocation"
		expectedResult = "all good"
	)

	myTmpl := template.New("some template")
	req := httpGetRequest(fmt.Sprintf("?location=%s", queryLocation))
	rr := httptest.NewRecorder()

	conditions := weather.Conditions{
		Location:    "API resolved location",
		Description: "description",
		Celsius:     5,
	}

	forecaster := forcasterMock{
		forecast: func(s string) (*weather.Conditions, error) {
			if s != queryLocation {
				t.Errorf("Unexpected argument in call to Forecast. Wanted %s but got %s", queryLocation, s)
			}
			return &conditions, nil
		},
	}
	rdr := &rendererMock{
		pathFunc: func(layouts ...string) []string { return layouts },
		buildFunc: func(layouts ...string) (*template.Template, error) {
			return myTmpl, nil
		},
		renderFunc: func(w http.ResponseWriter, tmpl *template.Template, v interface{}) error {
			if tmpl != myTmpl {
				t.Errorf("Unexpected argument in call to RenderTemplate. Wanted template '%s' got '%s'", myTmpl.Name(), tmpl.Name())
			}

			if m, ok := v.(map[string]interface{}); !ok {
				t.Error("Unexpected type in call to RenderTemplate. Want map[string]interface{} but got", reflect.TypeOf(v))
			} else {
				if err := checkMapFields(conditions, m); err != nil {
					t.Error(err)
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedResult))
			return nil
		},
	}

	http.HandlerFunc(widgetHandler(rdr, forecaster)).ServeHTTP(rr, req)

	if err := checkResponse(rr.Code, http.StatusOK,
		rr.Body.String(), expectedResult); err != nil {
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

		widgetHandler(&rendererMock{
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
				wanted: '%v' received '%v'`,
			expectedLayouts,
			layouts,
		)
	}
	return nil
}

func checkResponse(code, expectedCode int, body, expectedBody string) error {
	if code != expectedCode {
		return fmt.Errorf("handler returned wrong status code: Got %d want %d",
			code, expectedCode)
	}

	if body != expectedBody {
		return fmt.Errorf("handler returned unexpected body: Got '%s' want '%s'",
			body, expectedBody)
	}
	return nil
}

func checkMapFields(conditions weather.Conditions, m map[string]interface{}) error {
	expected := map[string]interface{}{
		"location":    conditions.Location,
		"celsius":     conditions.Celsius,
		"description": conditions.Description,
	}

	if !reflect.DeepEqual(expected, m) {
		return fmt.Errorf("unexpected arguments in call to RenderTemplate. Wanted '%v' but got '%v'", expected, m)
	}
	return nil
}
