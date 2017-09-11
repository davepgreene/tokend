package http

import (
	"encoding/json"
	"net/http"

	"github.com/davepgreene/tokend/api"
	"github.com/davepgreene/tokend/provider"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"io"
)

type kmsHandler struct {
	storage *api.Storage
}

type kmsPayload struct {
	Region     string
	Ciphertext string
	Datakey    string
}

func newKMSHandler(s *api.Storage) http.Handler {
	return handlers.MethodHandler{
		"POST": &kmsHandler{s},
	}
}

func (h *kmsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse the posted data and marshall it to a struct
	var body kmsPayload
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	switch {
	case err == io.EOF:
		body.Ciphertext = ""
		body.Datakey = ""
		body.Region = "us-east-1"
	case err != nil:
		panic(err)
	}

	ciphertext := body.Ciphertext
	datakey := body.Datakey
	region := body.Region

	if len(ciphertext) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}

	ch := make(chan provider.Provider)
	kms := provider.NewKMSProvider(&ch, map[string]string{
		ciphertext: ciphertext,
		datakey:    datakey,
		region:     region,
	})

	kms.Initialize()
	kmsResponse, err := kms.Get()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
	}

	data, _ := json.Marshal(kmsResponse)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
