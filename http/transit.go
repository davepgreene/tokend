package http

import (
	"encoding/json"
	"net/http"

	"github.com/davepgreene/tokend/api"
	"github.com/davepgreene/tokend/provider"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"sort"
)

type transitHandler struct {
	storage *api.Storage
}

type transitPayload struct {
	Ciphertext string `json:"ciphertext"`
	Key        string `json:"key"`
}

func newTransitHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"POST": &transitHandler{s},
	}
}

func (h *transitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	validFields := []string{"plaintext"}
	vars := mux.Vars(r)
	t := vars["token"]

	// Handle extracting the token from the vars map and passing it to the storage lookup

	// Parse the posted data and marshall it to a struct
	var body transitPayload
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	switch {
	case err == io.EOF:
		body.Ciphertext = ""
		body.Key = ""
	case err != nil:
		panic(err)
	}

	key := body.Key
	ciphertext := body.Ciphertext

	if len(key) == 0 || len(ciphertext) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}

	secret := "default"
	tokenManager := h.storage.Lookup(t, secret, "token", map[string]string{
		"token":  t,
		"secret": secret,
	})
	token, err := tokenManager.Initialize().Get()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}

	ch := make(chan provider.Provider)
	transit := provider.NewTransitProvider(&ch, map[string]string{
		"key":        key,
		"ciphertext": ciphertext,
		"token":      token["token"],
	})
	transit.Initialize()
	transitResponse, err := transit.Get()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
	}

	// Strip out anything but the token
	sort.Strings(validFields)
	for k := range transitResponse {
		i := sort.SearchStrings(validFields, k)
		if i > len(transitResponse) || validFields[i] != k {
			delete(transitResponse, k)
		}
	}
	data, _ := json.Marshal(transitResponse)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
