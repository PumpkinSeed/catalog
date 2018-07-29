package catalog

// @TODO setup healthcheck chain
func healthcheck(services map[identifier]*ServiceSpec) error {

	for _, service := range services {
		//go func() {}()
		// @TODO put it all into goroutines, channel if err
		if service.Healthcheck {
			alive, err := service.HealthcheckFunc()
			if err != nil {
				return err
			}

			service.IsAlive = alive
		}
	}

	return nil
}
