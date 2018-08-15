package api

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

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

		closeChan chan bool
	}{
		{
			name: "webserver",
			host: "localhost",
			port: 8080,
			tags: []string{"web", "http"},

			closeChan: make(chan bool),
		},
		{
			name: "webserver2",
			host: "localhost",
			port: 8081,
			tags: []string{"web", "http"},

			closeChan: make(chan bool),
		},
		{
			name: "auth",
			host: "localhost",
			port: 8001,
			tags: []string{"auth", "http"},

			closeChan: make(chan bool),
		},
	}
)

func init() {
	var server = catalog.NewServer(binAddr, hcStorage, &sync.RWMutex{})
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

func TestServices(t *testing.T) {
	services, err := testCatalogInstance.Services()
	if err != nil {
		t.Error(err)
	}

	for _, service := range services {
		if service.IsAlive != true {
			t.Errorf("%s should be alive", service.Name)
		}
	}

	testServices[1].closeChan <- true
	time.Sleep(4 * time.Second)
}

func TestFalseAliveServices(t *testing.T) {
	service, err := testCatalogInstance.Service(nil, &testServices[1].name)
	if err != nil {
		t.Error(err)
	}

	if service.IsAlive != false {
		t.Error("Service shouldn't be alive.")
	}
}

func TestTrueAliveServices(t *testing.T) {
	service, err := testCatalogInstance.Service(&testServices[0].id, nil)
	if err != nil {
		t.Error(err)
	}

	if service.IsAlive != true {
		t.Error("Service should be alive.")
	}
}

func TestDeregister(t *testing.T) {
	err := testCatalogInstance.Deregister(nil, &testServices[1].name)
	if err != nil {
		t.Error(err)
	}
}

func TestServicesAfterDeregister(t *testing.T) {
	services, err := testCatalogInstance.Services()
	if err != nil {
		t.Error(err)
	}

	if len(services) != 2 {
		t.Error("There should be 2 service in the storage no more")
	}
}

func startServices() {
	for _, service := range testServices {
		go func(addr string, closeChan chan bool) {
			l, err := net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}

			defer l.Close()
			for {
				select {
				case <-closeChan:
					fmt.Println("released")
					return
				default:
					conn, err := l.Accept()
					if err != nil {
						panic(err)
					}
					conn.Close()
				}
			}
		}(service.host+":"+strconv.Itoa(service.port), service.closeChan)
	}
}

func hcStorage(name string) (time.Duration, func() (bool, error)) {
	switch name {
	case "webserver":
		return 2 * time.Second, getCommonHCFunc(testServices[0].host + ":" + strconv.Itoa(testServices[0].port))
	case "webserver2":
		return 2 * time.Second, getCommonHCFunc(testServices[1].host + ":" + strconv.Itoa(testServices[1].port))
	case "auth":
		return 2 * time.Second, getCommonHCFunc(testServices[2].host + ":" + strconv.Itoa(testServices[2].port))
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
