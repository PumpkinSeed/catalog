package catalog

import (
	"time"
)

type Storage interface {
	Register(address string, port int, tags []string, additional interface{}) (identifier, error)
	Deregister(id identifier) error
	Service(id identifier) (ServiceSpec, error)
	Services() map[identifier]ServiceSpec
	SetupHealthcheck(id identifier, period time.Duration, f func() error) error
	Healthcheck(id identifier) error
}

// ServiceSpec represent the specification of a service
type ServiceSpec struct {
	ID      identifier
	Address string
	Port    int
	Tags    []string

	Healthcheck       bool
	HealthcheckFunc   func() error
	HealthcheckPeriod time.Duration
	IsAlive           bool

	Additional interface{}
}

type storage struct {
	services map[identifier]ServiceSpec
}

func NewStorage() Storage {
	return &storage{
		services: make(map[identifier]ServiceSpec),
	}
}

func (s *storage) Register(address string, port int, tags []string, additional interface{}) (identifier, error) {
	id := NewID()
	service := ServiceSpec{
		ID:         id,
		Address:    address,
		Port:       port,
		Tags:       tags,
		Additional: additional,
	}
	s.services[id] = service
	return id, nil
}
func (s *storage) Deregister(id identifier) error {
	delete(s.services, id)
	return nil
}
func (s *storage) Service(id identifier) (ServiceSpec, error) {
	service, ok := s.services[id]
	if !ok {
		return ServiceSpec{}, ErrUndefinedService
	}
	return service, nil
}
func (s *storage) Services() map[identifier]ServiceSpec {
	return s.services
}
func (s *storage) SetupHealthcheck(id identifier, period time.Duration, f func() error) error {
	// Check service before setup healthcheck
	err := f()
	if err != nil {
		return err
	}

	service, ok := s.services[id]
	if !ok {
		return ErrUndefinedService
	}

	service.Healthcheck = true
	service.HealthcheckPeriod = period
	service.HealthcheckFunc = f

	s.services[id] = service

	return nil
}
func (s *storage) Healthcheck(id identifier) error {
	if s.services[id].Healthcheck {
		return s.services[id].HealthcheckFunc()
	}
	return nil
}
