package catalog

// @TODO setup healthcheck chain
func healthcheck(services *map[identifier]ServiceSpec) error {

	for _, service := range *services {
		if service.Healthcheck {

		}
	}

	return nil
}
