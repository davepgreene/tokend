package http

import (
	"encoding/json"
	"net/http"

	"github.com/davepgreene/tokend/api"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

type tokenHandler struct {
	storage *api.Storage
}

func newTokenHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"GET": &tokenHandler{s},
	}
}

func (h *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := "default"
	secret := "default"
	tokenManager := h.storage.Lookup(t, secret, "token", map[string]string{
		"token":  t,
		"secret": secret,
	})
	tokenManager.Initialize()
	token, err := tokenManager.Get()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
	}

	data, _ := json.Marshal(token)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
