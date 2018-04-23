package tpl_test

import (
	"reflect"
	"testing"

	"github.com/wwgberlin/go-weather-widget/tpl"
)

func TestRenderer_PathToTemplateFiles(t *testing.T) {
	rdr := tpl.NewRenderer("some_path")
	expected := []string{"some_path/a", "some_path/b"}
	returned := rdr.PathToTemplateFiles("a", "b")
	if !reflect.DeepEqual(expected, returned) {
		t.Error("PathToTemplateFiles returned incorrect result. Wanted %v but got %v", expected, returned)
	}
}

func TestRenderer_BuildTemplate(t *testing.T) {
	//dir, _ := os.Getwd()
	//rdr := tpl.NewRenderer(filepath.Join(dir, "tpl"))
	//rdr := Renderer{}
	//build := rdr.BuildTemplate("some_name", "some", "files")
	//rdr.BuildTemplate()
}
