package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
)

type cubbyHoleHandler struct{}

func newCubbyHoleHandler() http.Handler {
	return handlers.MethodHandler{
		"GET": &cubbyHoleHandler{},
	}
}

func (h *cubbyHoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
