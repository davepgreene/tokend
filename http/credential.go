package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
)

type credentialHandler struct{}

func newCredentialHandler() http.Handler {
	return handlers.MethodHandler{
		"GET": &credentialHandler{},
	}
}

func (h *credentialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
