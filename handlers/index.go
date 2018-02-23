package handlers

import (
	"net/http"

	"github.com/ecosia/women-who-go/tpl"
)

// NewIndexHandler returns a new handler that is supposed
// to serve the app's index page
func NewIndexHandler() http.Handler {
	return http.HandlerFunc(handleIndex)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := tpl.Render(w, "index", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
