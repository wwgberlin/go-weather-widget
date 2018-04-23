package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var path = "./examples"

// basic templates: pass arguments to a template and execute
// http://127.0.0.1:8080/example1?name=Women%20Who%20Go
func example1(w http.ResponseWriter, r *http.Request) {
	//we typically pass in a map object
	params := map[string]interface{}{"Name": r.URL.Query().Get("name")}

	templateFiles := []string{
		filepath.Join(path, "example1", "hello.tmpl"),
	}

	if t, err := template.ParseFiles(templateFiles...); err == nil {
		if err := t.Execute(w, params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// nesting templates: defining multiple templates and executing them
// http://127.0.0.1:8080/example2?name=Women%20Who%20Go
func example2(w http.ResponseWriter, r *http.Request) {
	params := map[string]interface{}{"Name": r.URL.Query().Get("name")}

	templateFiles := []string{
		filepath.Join(path, "example2", "hello.tmpl"),
		filepath.Join(path, "example2", "head.tmpl"),
	}

	if t, err := template.ParseFiles(templateFiles...); err == nil {
		if err := t.Execute(w, params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// overriding template definition
func example3(w http.ResponseWriter, r *http.Request) {
	params := map[string]interface{}{"Name": r.URL.Query().Get("name")}

	templateFiles := []string{
		filepath.Join(path, "example3", "hello.tmpl"),
		filepath.Join(path, "example3", "head.tmpl"),
	}

	if t, err := template.ParseFiles(templateFiles...); err == nil {
		if err := t.Execute(w, params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// executing templates by name
// this method will produce and empty body - but check out the source
func example4(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query().Get("name")

	templateFiles := []string{
		filepath.Join(path, "example4", "hello.tmpl"),
		filepath.Join(path, "example4", "head.tmpl"),
	}

	if t, err := template.ParseFiles(templateFiles...); err == nil {
		if err := t.ExecuteTemplate(w, "nested", params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func foo(arg string) string {
	return "foo called with " + arg
}
func mkSlice(args ...int) []int {
	return args
}

func reverse(args []int) []int {
	for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
		args[i], args[j] = args[j], args[i]
	}
	return args
}

// Actual coding
func example5(w http.ResponseWriter, r *http.Request) {
	params := map[string]interface{}{"MyRange": []string{"1", "2", "3"}}

	templateFiles := []string{
		filepath.Join(path, "example5", "fun_stuff.tmpl"),
	}

	// Template provides a Must interface
	// When we define a named template, it will look it up.
	template.Must(template.New("fun").
		Funcs(template.FuncMap{
			"foo":     foo,
			"mkSlice": mkSlice,
			"reverse": reverse,
		}).ParseFiles(templateFiles...)).Execute(w, params)
}
func main() {
	// basic templates: pass arguments to a template and execute
	// http://127.0.0.1:8080/example1?name=Women%20Who%20Go
	http.HandleFunc("/example1", example1)

	// nesting templates: defining multiple templates and executing them
	// http://127.0.0.1:8080/example2?name=Women%20Who%20Go
	http.HandleFunc("/example2", example2)

	// overriding template definition
	http.HandleFunc("/example3", example3)

	// executing templates by name
	http.HandleFunc("/example4", example4)

	// actual go templates coding
	http.HandleFunc("/example5", example5)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
