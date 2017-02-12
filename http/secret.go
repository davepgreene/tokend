package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
)

type secretHandler struct{}

func newSecretHandler() http.Handler {
	return handlers.MethodHandler{
		"GET": &secretHandler{},
	}
}

func (h *secretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
