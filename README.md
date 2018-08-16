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

The Healthcheck storage is a `key => function` storage, provide manageable storage for the healthchecks of the services.

```
// Pass the hcStorage, it gets a name and returns the healthcheck function
// the mutex used for healthcheck thread-safety
var server = catalog.NewServer(binAddr, hcStorage, &sync.RWMutex{})

func hcStorage(name string) (time.Duration, func() (bool, error)) {
	switch name {
	case "webserver":
		return 2 * time.Second, getCommonHCFunc("127.0.0.1:8008")
	case "webserver2":
		return 2 * time.Second, getCommonHCFunc("127.0.0.1:8003")
	case "auth":
		return 2 * time.Second, getCommonHCFunc("127.0.0.1:8018")
	}

	return 2 * time.Second, nil
}

func getCommonHCFunc(address string) func() (bool, error) {
	return func() (bool, error) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return false, nil
		}
		defer conn.Close()

		err = conn.SetDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			return false, nil
		}
		//fmt.Println(address)

		return true, nil
	}
}
```
