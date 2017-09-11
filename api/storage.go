package api

import "fmt"

type Storage struct {
	timeout  int
	managers map[string]*LeaseManager
}

func NewStorage(timeout int) *Storage {
	managers := make(map[string]*LeaseManager)

	return &Storage{
		timeout:  timeout,
		managers: managers,
	}
}

func (s *Storage) Lookup(token string, secret string, provider string, config map[string]string) *LeaseManager {
	manager := s.getLeaseManager(token, secret, provider, config)
	return manager
}

func (s *Storage) getLeaseManager(token string, secret string, provider string, config map[string]string) *LeaseManager {
	name := fmt.Sprintf("%s/%s/%s", provider, token, secret)
	if manager, ok := s.managers[name]; ok {
		return manager
	}

	manager := NewLeaseManager(provider, name, config)

	s.managers[name] = manager
	return manager
}
