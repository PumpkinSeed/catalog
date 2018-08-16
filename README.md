# catalog

First of all it's a single-node in-memory service discovery tool for test behaviours depends on small services. On the other hand it implements Consul catalog API as an in-memory mock.

#### Usage in Go test cases

```
// ----- start the mock as a goroutine
var binAddr = "127.0.0.1:7777"
var server = catalog.NewServer(binAddr, nil, &sync.RWMutex{})
go func() {
	panic(server.Listen())
}()

// ----- create a new API instance
catalogInstance = api.NewCatalog(binAddr)

// ----- register service
var (
	nameOfService = "webserver"
	hostOfService = "localhost"
	portOfService = 8080
	tagsOfService = []string{"web", "http"}
)
idOfService, err := catalogInstance.Register(nameOfService, hostOfService, portOfService, tagsOfService, nil)
// handler err

// ----- get service by id
service, err := catalogInstance.Service(&idOfService, nil)
// handler err

// ----- get service by name
service, err := catalogInstance.Service(nil, &nameOfService)
// handler err

// ----- get all services
services, err := catalogInstance.Services()
// handler err

// ----- deregister service by name
err := catalogInstance.Deregister(nil, &nameOfService)
// handler err
```

#### Usage for different languages

There is only Go API implementation for the catalog, but feel free the compile it to a binary and call the TCP socket endpoints. `Standalone directory`

#### Healthcheck storage

Later...
