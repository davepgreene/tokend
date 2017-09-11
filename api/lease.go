package api

import (
	"errors"
	p "github.com/davepgreene/tokend/provider"
	"github.com/davepgreene/tokend/utils"
	"github.com/davepgreene/tokend/utils/status"
	"github.com/satori/go.uuid"
)

type LeaseManager struct {
	status        status.Status
	data          map[string]string
	leaseDuration int
	provider      p.Provider
	name          string
	error         error
	correlationID uuid.UUID
	channel       chan p.Provider
}

func NewLeaseManager(provider string, name string, conf map[string]string) *LeaseManager {
	c := make(chan p.Provider)
	prov, _ := p.CreateProvider(provider, &c, conf)

	return &LeaseManager{
		correlationID: utils.CreateCorrelationID(),
		provider:      prov,
		status:        status.Pending,
		name:          name,
		channel:       c,
	}
}

func (l *LeaseManager) Renewable() bool {
	return l.provider.Renewable()
}

func (l *LeaseManager) Expires() bool {
	return l.provider.Expires()
}

func (l *LeaseManager) Initialize() *LeaseManager {
	err := l.provider.Initialize()
	if err != nil {
		panic(err)
	}

	data, err := l.provider.Get()
	if err != nil {
		panic(err)
	}
	l.data = data
	l.leaseDuration = l.provider.LeaseDuration()
	l.status = status.Ready

	return l.Renew()
}

func (l *LeaseManager) Renew() *LeaseManager {
	return l
}

func (l *LeaseManager) Get() (map[string]string, error) {
	if l.status == status.Ready {
		return l.provider.Get()
	}

	return nil, errors.New("LeaseManager isn't ready yet")
}
