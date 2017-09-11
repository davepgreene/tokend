package transport

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Warden is a container for data returned
// from a request for a token from Warden
type Warden struct {
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreationTime   time.Time `json:"creation_time"`
	LeaseDuration  int       `json:"lease_duration"`
	Renewable      bool      `json:"renewable"`
	metadata       Metadata
}

type wardenTransport struct {
	ClientToken    string            `json:"client_token"`
	Policies       []string          `json:"policies"`
	Metadata       map[string]string `json:"metadata"`
	LeaseDuration  int               `json:"lease_duration"`
	Renewable      bool              `json:"renewable"`
	CreationTime   string            `json:"creation_time"`
	ExpirationTime string            `json:"expiration_time"`
}

// NewWarden creates a new DTO for Warden data
func NewWarden(m *Metadata) *Warden {
	w := Warden{
		metadata: *m,
	}
	return w.Send()
}

// Send POSTs the Metadata documents to Warden
// MUTABLE
func (w *Warden) Send() *Warden {
	transport := send(&w.metadata)

	expiration, _ := time.Parse(time.RFC3339Nano, transport.ExpirationTime)
	creation, _ := time.Parse(time.RFC3339Nano, transport.CreationTime)

	w.Token = transport.ClientToken
	w.ExpirationTime = expiration
	w.CreationTime = creation
	w.LeaseDuration = transport.LeaseDuration
	w.Renewable = transport.Renewable

	return w
}

func send(m *Metadata) *wardenTransport {
	endpoint := fmt.Sprintf("http://%s:%d%s",
		viper.GetString("warden.host"),
		viper.GetInt("warden.port"),
		viper.GetString("warden.path"))
	body := HTTPPost(endpoint, make(map[string]string), &m)

	t := wardenTransport{}
	if err := json.Unmarshal(body, &t); err != nil {
		panic(err)
	}

	return &t
}
