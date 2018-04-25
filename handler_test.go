package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/wwgberlin/go-weather-widget/weather"
)

type (
	rendererMock struct {
		buildInvoked  bool
		renderInvoked bool
		buildFunc     func(...string) *template.Template
		renderFunc    func(io.Writer, *template.Template, interface{}) error
	}
	forecasterMock struct {
		invoked  bool
		forecast func(string) (*weather.Conditions, error)
	}
)

func (f forecasterMock) Forecast(s string) (*weather.Conditions, error) {
	return f.forecast(s)
}

func (rdr *rendererMock) BuildTemplate(dep ...string) *template.Template {
	rdr.buildInvoked = true
	return rdr.buildFunc(dep...)
}

func (rdr *rendererMock) RenderTemplate(w io.Writer, tmpl *template.Template, data interface{}) error {
	rdr.renderInvoked = true
	return rdr.renderFunc(w, tmpl, data)
}

func TestIndexHandler_BuildTemplate(t *testing.T) {
	expectedFiles := []string{
		"my/path/index.tmpl",
		"my/path/layouts/layout.tmpl",
		"my/path/layouts/head.tmpl",
	}

	rdr := &rendererMock{
		buildFunc: func(layouts ...string) *template.Template {
			if err := checkTemplates(layouts, expectedFiles); err != nil {
				t.Error(err)
			}
			return template.New("some template")
		},
	}

	indexHandler("my/path/", rdr)

	if !rdr.buildInvoked {
		t.Error("BuildTemplated was expected to be called")
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
		buildFunc: func(layouts ...string) *template.Template {
			return myTmpl
		},
		renderFunc: func(w io.Writer, tmpl *template.Template, i interface{}) error {
			if tmpl != myTmpl {
				t.Errorf("Unexpected argument in call to RenderTemplate. Wanted template %s got %s", myTmpl.Name(), tmpl.Name())
			}
			w.Write([]byte(expectedResult))
			return nil
		},
	}

	http.HandlerFunc(indexHandler("", rdr)).ServeHTTP(rr, req)

	if err := checkResponse(rr.Code, http.StatusOK,
		rr.Body.String(), expectedResult); err != nil {
		t.Error(strings.Title(err.Error()))
	}
}

func TestIndexHandler_FailToRenderTemplate(t *testing.T) {
	const errMsg = "some error occurred"

	req := httpGetRequest("some_path")
	rr := httptest.NewRecorder()
	rdr := &rendererMock{
		buildFunc: func(s2 ...string) *template.Template {
			return nil
		},
		renderFunc: func(io.Writer, *template.Template, interface{}) error {
			return errors.New(errMsg)
		},
	}

	http.HandlerFunc(indexHandler("", rdr)).ServeHTTP(rr, req)
	body := rr.Body.String()[0 : len(rr.Body.String())-1]

	if err := checkResponse(
		rr.Code,
		http.StatusInternalServerError,
		body,
		errMsg,
	); err != nil {
		t.Error(strings.Title(err.Error()))
	}
}

func TestWidgetHandler_TestBuild(t *testing.T) {
	expectedFiles := []string{
		"my/path/widget.tmpl",
		"my/path/layouts/layout.tmpl",
		"my/path/layouts/head.tmpl",
	}

	rdr := &rendererMock{
		buildFunc: func(layouts ...string) *template.Template {
			if err := checkTemplates(layouts, expectedFiles); err != nil {
				t.Error(err)
			}
			return template.New("some template")
		},
	}

	widgetHandler("my/path/", rdr, forecasterMock{})

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

	forecaster := forecasterMock{
		forecast: func(s string) (*weather.Conditions, error) {
			if s != queryLocation {
				t.Errorf("Unexpected argument in call to Forecast. Wanted %s but got %s", queryLocation, s)
			}
			return &conditions, nil
		},
	}
	rdr := &rendererMock{
		buildFunc: func(layouts ...string) *template.Template {
			return myTmpl
		},
		renderFunc: func(w io.Writer, tmpl *template.Template, v interface{}) error {
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
			w.Write([]byte(expectedResult))
			return nil
		},
	}

	http.HandlerFunc(widgetHandler("", rdr, forecaster)).ServeHTTP(rr, req)

	if err := checkResponse(rr.Code, http.StatusOK,
		rr.Body.String(), expectedResult); err != nil {
		t.Error(strings.Title(err.Error()))
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
