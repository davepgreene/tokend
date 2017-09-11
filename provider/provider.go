package provider

import "time"

// Provider interface
type Provider interface {
	Initialize() error
	Renew() error
	Invalidate()
	Get() (map[string]string, error)
	Renewable() bool
	Expires() bool
	LeaseDuration() int
}

// GenericProvider is an embeddable struct that contains
// common elements between Provider types
type GenericProvider struct {
	channel        *chan Provider
	token          string
	renewable      bool
	path           string
	expirationTime time.Time
	creationTime   time.Time
}

// Renewable determines whether the Provider can be renewed
func (t *GenericProvider) Renewable() bool {
	return t.renewable
}

// Expires checks whether the Provider can expire
func (t *GenericProvider) Expires() bool {
	expiration := t.expirationTime
	creation := t.creationTime
	if expiration.IsZero() || creation.IsZero() {
		return false
	}

	return true
}

func (t *GenericProvider) LeaseDuration() int {
	return 0
}

type providerData struct {
	Data map[string]string `json:"data"`
}

func (t *providerData) Get() map[string]string {
	return t.Data
}
