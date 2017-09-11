package provider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/davepgreene/tokend/transport"
)

type TransitProvider struct {
	GenericProvider
	data       transitProviderData
	key        string
	ciphertext string
}

type transitProviderData struct {
	providerData
	transitData
}

type transitPayload struct {
	Ciphertext string `json:"ciphertext"`
}

type transitData struct {
	Plaintext string `json:"plaintext"`
}

type transitResponse struct {
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          transitData `json:"data"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      interface{} `json:"warnings"`
	Auth          interface{} `json:"auth"`
}

func NewTransitProvider(ch *chan Provider, params map[string]string) Provider {
	return &TransitProvider{
		key:        params["key"],
		ciphertext: params["ciphertext"],
		GenericProvider: GenericProvider{
			token:     params["token"],
			renewable: false,
			channel:   ch,
		},
	}
}

func (t *TransitProvider) Initialize() error {
	body := transitPayload{
		Ciphertext: t.ciphertext,
	}
	vault := vaultConfig()
	endpoint := fmt.Sprintf("%s/v1/transit/decrypt/%s",
		vault.Address,
		t.key)

	headers := make(map[string]string)
	headers["X-Vault-Token"] = t.token
	resp := transport.HTTPPost(endpoint, headers, body)

	r := transitResponse{}
	err := json.Unmarshal(resp, &r)
	if err != nil {
		panic(err)
	}

	// Decode the plaintext response
	plaintext, err := base64.StdEncoding.DecodeString(r.Data.Plaintext)
	if err != nil {
		panic(err)
	}

	t.data = transitProviderData{
		transitData: transitData{
			Plaintext: string(plaintext),
		},
	}

	return nil
}

func (t *TransitProvider) Renew() error {
	return t.Initialize()
}

// Invalidate is a stub to conform to the ProviderInterface
func (t *TransitProvider) Invalidate() {}

func (t *TransitProvider) Get() (map[string]string, error) {
	return map[string]string{
		"plaintext": t.data.Plaintext,
	}, nil
}
