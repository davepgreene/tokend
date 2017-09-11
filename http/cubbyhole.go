package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/davepgreene/tokend/api"
)

type cubbyHoleHandler struct{
	storage *api.Storage
}

func newCubbyHoleHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"GET": &cubbyHoleHandler{s},
	}
}

func (h *cubbyHoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotImplemented)
	b, _ := json.Marshal(h)

	w.Write(b)
}
