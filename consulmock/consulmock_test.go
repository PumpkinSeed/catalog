package consulmock

import (
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/PumpkinSeed/catalog"
	"github.com/PumpkinSeed/consul/api"
)

var (
	testCatalogInstance Catalog
	binAddr             = "127.0.0.1:8889"
	testServices        = []*struct {
		name       string
		host       string
		port       int
		tags       []string
		registered bool
		isAlive    bool // for later
	}{
		{
			name: "websrever",
			host: "localhost",
			port: 8080,
			tags: []string{"web", "http"},
		},
		{
			name: "websrever2",
			host: "localhost",
			port: 8081,
			tags: []string{"web", "http"},
		},
		{
			name: "auth",
			host: "localhost",
			port: 8001,
			tags: []string{"auth", "http"},
		},
	}
)

func init() {
	var server = catalog.NewServer(binAddr, nil, &sync.RWMutex{})
	go func() {
		panic(server.Listen())
	}()
	startServices()
}

func TestNewCatalog(t *testing.T) {
	testCatalogInstance = NewCatalog(binAddr)
}

func TestRegister(t *testing.T) {
	for _, service := range testServices {
		testCatalogInstance.Register(&api.CatalogRegistration{
			Service: &api.AgentService{
				ID:      service.name,
				Address: service.host,
				Port:    service.port,
				Tags:    service.tags,
			},
		}, nil)
	}
}

func startServices() {
	for _, service := range testServices {
		go func(addr string) {
			l, err := net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}

			defer l.Close()
			for {
				conn, err := l.Accept()
				if err != nil {
					panic(err)
				}
				conn.Close()
			}
		}(service.host + ":" + strconv.Itoa(service.port))
	}
}
