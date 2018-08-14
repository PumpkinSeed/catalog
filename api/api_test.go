package api

import (
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/PumpkinSeed/catalog"
)

var (
	testCatalogInstance Catalog
	binAddr             = "127.0.0.1:7777"
	testServices        = []*struct {
		id         string
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
		id, err := testCatalogInstance.Register(service.name, service.host, service.port, service.tags, nil)
		if err != nil {
			t.Error(err)
		}

		service.id = id
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
