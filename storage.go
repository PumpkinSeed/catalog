package catalog

import (
	"strconv"
	"sync"
	"time"
)

type Storage interface {
	Register(name string, host string, port int, tags []string, additional interface{}) (identifier, error)
	Deregister(id *identifier, name *string) error
	Service(id *identifier, name *string) (*ServiceSpec, error)
	Services() map[identifier]*ServiceSpec
	SetupHealthcheck(id identifier, period time.Duration, f func() (bool, error)) error
	Healthcheck() error
}

// ServiceSpec represent the specification of a service
type ServiceSpec struct {
	ID      identifier `json:"id"`
	Name    string     `json:"name"`
	Host    string     `json:"host"`
	Port    int        `json:"port"`
	Address string     `json:"address"`
	Tags    []string   `json:"tags"`

	Healthcheck       bool                 `json:"healthcheck"`
	HealthcheckFunc   func() (bool, error) `json:"-"`
	HealthcheckPeriod time.Duration        `json:"healthcheck_period"`
	IsAlive           bool                 `json:"is_alive"`

	Additional interface{}
}

type storage struct {
	mutex              sync.RWMutex
	services           map[identifier]*ServiceSpec
	healthcheckStorage func(name string) (time.Duration, func() (bool, error))
	healthcheckMutex   sync.RWMutex
}

func NewStorage(healthcheckStorage func(name string) (time.Duration, func() (bool, error)), mutex sync.RWMutex) Storage {
	return &storage{
		services:           make(map[identifier]*ServiceSpec),
		healthcheckStorage: healthcheckStorage,
		healthcheckMutex:   sync.RWMutex{},
		mutex:              mutex,
	}
}

func (s *storage) Register(name string, host string, port int, tags []string, additional interface{}) (identifier, error) {
	var hcFunc func() (bool, error)
	var period time.Duration
	if s.healthcheckStorage != nil {
		period, hcFunc = s.healthcheckStorage(name)
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := NewID()
	service := ServiceSpec{
		ID:         id,
		Name:       name,
		Host:       host,
		Port:       port,
		Address:    host + ":" + strconv.Itoa(port),
		Tags:       tags,
		Additional: additional,
	}
	s.services[id] = &service

	s.SetupHealthcheck(id, period, hcFunc)
	return id, nil
}

func (s *storage) Deregister(id *identifier, name *string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// ID first manner
	if id != nil {
		delete(s.services, *id)
		return nil
	} else if name != nil {
		ss := s.findByName(*name)
		delete(s.services, ss.ID)
		return nil
	}

	return ErrServiceRequestInvalid
}

func (s *storage) Service(id *identifier, name *string) (*ServiceSpec, error) {
	var service *ServiceSpec
	var ok bool

	s.mutex.RLock()

	// ID first manner
	if id != nil {
		service, ok = s.services[*id]
	} else if name != nil {
		ss := s.findByName(*name)
		service, ok = s.services[ss.ID]
	} else {
		return service, ErrServiceRequestInvalid
	}
	s.mutex.RUnlock()
	if !ok {
		return &ServiceSpec{}, ErrUndefinedService
	}
	return service, nil
}

func (s *storage) Services() map[identifier]*ServiceSpec {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.services
}

func (s *storage) SetupHealthcheck(id identifier, period time.Duration, f func() (bool, error)) error {
	// Check service before setup healthcheck
	if f == nil {
		return nil
	}
	alive, err := f()
	if err != nil {
		return err
	}

	service, ok := s.services[id]
	if !ok {
		return ErrUndefinedService
	}
	if alive {
		service.IsAlive = true
	}

	service.Healthcheck = true
	service.HealthcheckPeriod = period
	service.HealthcheckFunc = f

	s.services[id] = service

	return nil
}

func (s *storage) Healthcheck() error {
	return healthcheck(s.services, s.healthcheckMutex)
}

func (s *storage) findByName(name string) *ServiceSpec {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, service := range s.services {
		if service.Name == name {
			var ss *ServiceSpec
			ss = service
			return ss
		}
	}

	return nil
}
