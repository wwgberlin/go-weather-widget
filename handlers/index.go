package handlers

import (
	"net/http"

	"github.com/ecosia/women-who-go/tpl"
)

// NewIndexHandler returns a new handle that is supposed
// to serve the app's index page
func NewIndexHandler() http.Handler {
	return http.HandlerFunc(handleIndex)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tpl.Render(w, "index", nil)
}
