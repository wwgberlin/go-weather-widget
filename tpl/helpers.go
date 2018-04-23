package tpl

import (
	"html/template"
	"regexp"
	"strings"
)

var helpers = template.FuncMap{
	"title":     strings.Title,
	"clothings": clothings,
}

var umbrellaStr = regexp.MustCompile("[[R|r]ain|[D|d]rizzl|[S|s]leet")

func clothings(weatherDesc string, celsius int) (clothes []string) {
	if umbrellaStr.Match([]byte(weatherDesc)) {
		clothes = append(clothes, "umbrella")
	}
	if celsius > 15 {
		clothes = append(clothes, "tshirt")
	}
	if celsius > 20 {
		clothes = append(clothes, "sunglasses")
	}
	if celsius > 22 {
		clothes = append(clothes, "hat")
	}
	if celsius < 15 {
		clothes = append(clothes, "boots")
		clothes = append(clothes, "scarf")
	}
	if celsius <= 15 {
		clothes = append(clothes, "coat")
	}
	return clothes
}
