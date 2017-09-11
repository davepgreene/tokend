package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/davepgreene/tokend/api"
)

type credentialHandler struct{
	storage *api.Storage
}

func newCredentialHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"GET": &credentialHandler{s},
	}
}

func (h *credentialHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
