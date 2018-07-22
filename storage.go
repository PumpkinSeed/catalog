package catalog

import (
	"strconv"
	"sync"
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
	ID      identifier `json:"id"`
	Host    string     `json:"host"`
	Port    int        `json:"port"`
	Address string     `json:"address"`
	Tags    []string   `json:"tags"`

	Healthcheck       bool          `json:"healthcheck"`
	HealthcheckFunc   func() error  `json:"-"`
	HealthcheckPeriod time.Duration `json:"healthcheck_period"`
	IsAlive           bool          `json:"is_alive"`

	Additional interface{}
}

type storage struct {
	sync.RWMutex
	services map[identifier]ServiceSpec
}

func NewStorage() Storage {
	return &storage{
		services: make(map[identifier]ServiceSpec),
	}
}

func (s *storage) Register(host string, port int, tags []string, additional interface{}) (identifier, error) {
	s.Lock()
	defer s.Unlock()
	id := NewID()
	service := ServiceSpec{
		ID:         id,
		Host:       host,
		Port:       port,
		Address:    host + ":" + strconv.Itoa(port),
		Tags:       tags,
		Additional: additional,
	}
	s.services[id] = service
	return id, nil
}

func (s *storage) Deregister(id identifier) error {
	s.Lock()
	defer s.Unlock()
	delete(s.services, id)
	return nil
}

func (s *storage) Service(id identifier) (ServiceSpec, error) {
	s.RLock()
	service, ok := s.services[id]
	s.RUnlock()
	if !ok {
		return ServiceSpec{}, ErrUndefinedService
	}
	return service, nil
}

func (s *storage) Services() map[identifier]ServiceSpec {
	s.RLock()
	defer s.RUnlock()
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
