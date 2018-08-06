package catalog

import (
	"sync"
)

// @TODO setup healthcheck chain
func healthcheck(services map[identifier]*ServiceSpec, mutex *sync.RWMutex) error {

	for _, service := range services {

		//go func() {}()
		// @TODO put it all into goroutines, channel if err
		if service.Healthcheck {
			alive, err := service.HealthcheckFunc()
			if err != nil {
				return err
			}

			mutex.Lock()
			service.IsAlive = alive
			mutex.Unlock()
		}
	}
	//mutex.Unlock()

	return nil
}
