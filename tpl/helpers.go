package tpl

import (
	"html/template"
	"regexp"
	"strings"
)

var Helpers = template.FuncMap{
	"title":  strings.Title,
	"clothe": clothe,
}

var umbrellaStr = regexp.MustCompile("[[R|r]ain|[D|d]rizzl|[S|s]leet")

func clothe(weatherDesc string, celsius int) (clothes []string) {
	if umbrellaStr.Match([]byte(weatherDesc)) {
		clothes = append(clothes, "umbrella")
	} else {
		if celsius > 22 {
			clothes = append(clothes, "hat")
		} else if celsius <= 10 {
			clothes = append(clothes, "winterhat")
		}
	}
	if celsius > 15 {
		clothes = append(clothes, "tshirt")
	}
	if celsius > 20 {
		clothes = append(clothes, "sunglasses")
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
