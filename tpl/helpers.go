package tpl

import (
	"html/template"
	"strings"
)

var helpers = template.FuncMap{
	"concat": concat,
	"title":  strings.Title,
}

func concat(tokens ...string) string {
	return strings.Join(tokens, "")
}
