package catalog

import (
	"sync"
)

// @TODO setup healthcheck chain
func healthcheck(services map[identifier]*ServiceSpec, mutex *sync.RWMutex) error {
	var errChan = make(chan error)
	var counter = 0
	for _, service := range services {
		counter++
		go func(service *ServiceSpec) {
			defer func() {
				counter--
			}()
			if service.Healthcheck {
				alive, err := service.HealthcheckFunc()
				if err != nil {
					errChan <- err
					return
				}

				mutex.Lock()
				service.IsAlive = alive
				mutex.Unlock()
			}
			return
		}(service)
	}

	go func() {
		for {
			if counter == 0 {
				errChan <- nil
				break
			}
		}
	}()

	for {
		select {
		case err := <-errChan:
			return err
		}
	}
}
