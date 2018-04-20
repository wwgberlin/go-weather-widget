package tpl

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

var helpers = template.FuncMap{
	"concat":     concat,
	"title":      strings.Title,
	"classNames": classNames,
	"clothings":  clothings,
}

func classNames(classes []string) string {
	fmt.Println(classes)
	var result []string
	for _, c := range classes {
		found := false
		for _, r := range result {
			if r == c {
				found = true
				break
			}
		}
		if !found {
			result = append(result, c)
		}
	}
	return strings.Join(result, " ")
}

var umbrellaStr = regexp.MustCompile("[[R|r]ain|[D|d]rizzl|[S|s]leet")

func clothings(weatherDesc string, celsius int) (clothes []string) {
	if umbrellaStr.Match([]byte(weatherDesc)) {
		clothes = append(clothes, "umbrella")
	}
	switch {
	case celsius > 12:
		clothes = append(clothes, "tshirt")
		fallthrough
	case celsius > 20:
		clothes = append(clothes, "sunglasses")
		fallthrough
	case celsius > 22:
		clothes = append(clothes, "hat")
	case celsius < 15:
		clothes = append(clothes, "boots")
		clothes = append(clothes, "scarf")
		fallthrough
	case celsius <= 12:
		clothes = append(clothes, "coat")
	}
	return clothes
}

func concat(tokens ...string) string {
	return strings.Join(tokens, "")
}
