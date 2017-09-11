package provider

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type factory func(ch *chan Provider, conf map[string]string) Provider

var registeredProviders = make(map[string]factory)

func Register(name string, factory factory) {
	if factory == nil {
		log.Panic(fmt.Sprintf("Provider factory %s does not exist.", name))
	}
	_, registered := registeredProviders[name]
	if registered {
		log.WithFields(log.Fields{
			name: name,
		}).Error(fmt.Sprintf("Provider factory %s already registered. Ignoring.", name))
	}
	registeredProviders[name] = factory
}

func CreateProvider(name string, ch *chan Provider, conf map[string]string) (Provider, error) {
	p, ok := registeredProviders[name]

	if !ok {
		// Factory has not been registered.
		// Make a list of all available provider factories for logging.
		availableProviders := make([]string, len(registeredProviders))
		for k := range registeredProviders {
			availableProviders = append(availableProviders, k)
		}
		return nil, fmt.Errorf("invalid Provider name. Must be one of: %s", strings.Join(availableProviders, ", "))
	}

	// Run the factory with the configuration.
	return p(ch, conf), nil
}

func init() {
	Register("token", NewTokenProvider)
	Register("transit", NewTransitProvider)
	Register("kms", NewKMSProvider)
}
