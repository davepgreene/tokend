package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
)

type healthHandler struct{}

func newHealthHandler() http.Handler {
	return handlers.MethodHandler{
		"GET": &healthHandler{},
	}
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
