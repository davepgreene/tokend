package transport

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var endpoints = map[string]string{
	"Document":  "/latest/dynamic/instance-identity/document",
	"Signature": "/latest/dynamic/instance-identity/signature",
	"Pkcs7":     "/latest/dynamic/instance-identity/pkcs7",
}

// Metadata is a type wrapper for EC2 Metadata operations
type Metadata struct {
	Document  string
	Signature string
	Pkcs7     string
}

// NewMetadata constructs a new Metadata type with a client
func NewMetadata() *Metadata {
	return &Metadata{}
}

func endpoint() string {
	return fmt.Sprintf("http://%s:%d", viper.GetString("metadata.host"), viper.GetInt("metadata.port"))
}

// Get retrieves documents from each important endpoint
func (m *Metadata) Get() *Metadata {
	for k, v := range endpoints {
		log.WithField("endpoint", endpoint())
		body := HTTPGet(endpoint()+v, make(map[string]string))

		field := reflect.ValueOf(m).Elem().FieldByName(k)

		if field.IsValid() && field.CanSet() {
			field.SetString(string(body))
		}
	}
	return m
}
