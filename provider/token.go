package provider

import (
	"time"

	"github.com/davepgreene/tokend/transport"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

type TokenProvider struct {
	GenericProvider
	expirationTime time.Time
	creationTime   time.Time
	warden         *transport.Warden
	metadata       *transport.Metadata
	data           tokenProviderData
}

type tokenProviderData struct {
	providerData
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
}

// NewTokenProvider instantiates and initializes a new TokenProvider
func NewTokenProvider(ch *chan Provider, params map[string]string) Provider {
	log.Info("Creating TokenProvider")
	return &TokenProvider{
		GenericProvider: GenericProvider{
			channel: ch,
		},
	}
}

// Initialize a token provider
func (t *TokenProvider) Initialize() error {
	if t.token != "" {
		return nil
	}

	t.metadata = t.getDocument()
	t.warden = t.sendDocument()

	t.token = t.warden.Token
	duration := t.warden.LeaseDuration

	t.data = tokenProviderData{
		LeaseID:       t.token,
		LeaseDuration: duration,
	}
	t.expirationTime = t.warden.ExpirationTime
	t.creationTime = t.warden.CreationTime
	t.renewable = true

	// TODO: REMOVE BEFORE USE
	log.WithFields(log.Fields{
		"data": t.data,
	}).Debug("Initialized data")

	return nil
}

// Renew sends a renewal request to Vault for the Provider's token
func (t *TokenProvider) Renew() error {
	if t.token == "" {
		t.Initialize()
	}

	// Renew actions against Vault
	client, err := api.NewClient(vaultConfig())
	if err != nil {
		return err
	}

	increment := viper.GetInt("vault.token_renew_increment")
	token, err := client.Auth().Token().Renew(t.token, increment)
	if err != nil {
		return err
	}
	t.token = token.LeaseID

	*t.channel <- t

	return nil
}

// Invalidate clears a Provider's underlying data
func (t *TokenProvider) Invalidate() {
	t.token = ""
	t.data = tokenProviderData{}
	t.expirationTime = time.Time{}
	t.creationTime = time.Time{}
}

// Get retrieves a Provider's data
func (t *TokenProvider) Get() (map[string]string, error) {
	return map[string]string{
		"token": t.data.LeaseID,
		//"lease_duration": strconv.Itoa(int(t.data.LeaseDuration)),
	}, nil
	//return t.data.Get(), nil
}

// CreationTime gets a Provider's creationTime
func (t *TokenProvider) CreationTime() time.Time {
	return t.creationTime
}

// ExpirationTime gets a Provider's expirationTime
func (t *TokenProvider) ExpirationTime() time.Time {
	return t.expirationTime
}

func (t *TokenProvider) LeaseDuration() int {
	if t.token == "" {
		t.Initialize()
	}

	return t.data.LeaseDuration
}

func (t *TokenProvider) getDocument() *transport.Metadata {
	if t.metadata != nil {
		return t.metadata
	}

	return transport.NewMetadata().Get()
}

func (t *TokenProvider) sendDocument() *transport.Warden {
	if t.warden != nil {
		return t.warden
	}

	if t.metadata == nil {
		t.metadata = t.getDocument()
	}

	return transport.NewWarden(t.metadata)
}

func (t *tokenProviderData) Get() map[string]string {
	return map[string]string{
		"token":          t.LeaseID,
		"lease_duration": strconv.Itoa(t.LeaseDuration),
	}
}
