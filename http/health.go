package http

import (
	"encoding/json"
	"net/http"

	"github.com/davepgreene/tokend/api"
	"github.com/gorilla/handlers"
)

type healthHandler struct {
	storage *api.Storage
}

func newHealthHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"GET": &healthHandler{s},
	}
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
