package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/davepgreene/tokend/api"
)

type secretHandler struct{
	storage *api.Storage
}

func newSecretHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"GET": &secretHandler{s},
	}
}

func (h *secretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
