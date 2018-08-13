package catalog

import (
	"net"
	"sync"
	"testing"
	"time"
)

const localhost = "127.0.0.1"

var services = []struct {
	id    Identifier
	name  string
	host  string
	port  string
	alive bool
}{
	{
		id:    NewID(),
		name:  "webserver",
		host:  localhost,
		port:  "8091",
		alive: true,
	},
	{
		id:    NewID(),
		name:  "authservice",
		host:  localhost,
		port:  "7000",
		alive: true,
	},
	{
		id:    NewID(),
		name:  "searchservice",
		host:  localhost,
		port:  "8088",
		alive: true,
	},
	{
		id:    NewID(),
		name:  "mailservice",
		host:  localhost,
		port:  "8001",
		alive: false,
	},
}

var serviceSpecs = map[Identifier]*ServiceSpec{
	services[0].id: {
		ID:              services[0].id,
		Name:            "webserver",
		Host:            localhost,
		Port:            8091,
		Healthcheck:     true,
		HealthcheckFunc: getCommonHCFunc(localhost + ":8091"),
		IsAlive:         true,
	},
	services[1].id: {
		ID:              services[1].id,
		Name:            "authservice",
		Host:            localhost,
		Port:            7000,
		Healthcheck:     true,
		HealthcheckFunc: getCommonHCFunc(localhost + ":7000"),
		IsAlive:         true,
	},
	services[2].id: {
		ID:              services[2].id,
		Name:            "searchservice",
		Host:            localhost,
		Port:            8088,
		Healthcheck:     true,
		HealthcheckFunc: getCommonHCFunc(localhost + ":8088"),
		IsAlive:         true,
	},
	services[3].id: {
		ID:              services[3].id,
		Name:            "mailservice",
		Host:            localhost,
		Port:            8001,
		Healthcheck:     true,
		HealthcheckFunc: getCommonHCFunc(localhost + ":8001"),
		IsAlive:         true,
	},
}

func TestHealthcheck(t *testing.T) {
	startServices(t)

	mutex := sync.RWMutex{}
	err := healthcheck(serviceSpecs, &mutex)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(4 * time.Second)
	if serviceSpecs[services[3].id].IsAlive != false {
		t.Errorf("Service spec with id %v should have false IsAlive", services[3].id)
	}
}

func startServices(t *testing.T) {
	var mutex = &sync.Mutex{}

	for _, service := range services {
		if !service.alive {
			continue
		}

		go func(addr string) {

			l, err := net.Listen("tcp", addr)
			if err != nil {
				t.Errorf("Error listening: %s", err.Error())
				return
			}

			defer l.Close()
			//fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
			for {
				// Listen for an incoming connection.
				mutex.Lock()
				conn, err := l.Accept()
				mutex.Unlock()
				if err != nil {
					t.Errorf("Error accepting: %s", err.Error())
				}
				// Handle connections in a new goroutine.
				conn.Close()
			}
		}(service.host + ":" + service.port)
	}
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

		return true, nil
	}
}
