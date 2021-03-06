package catalog

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	var httpPort = 8082
	var id Identifier

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			message := r.URL.Path
			message = strings.TrimPrefix(message, "/")
			message = "Hello " + message
			w.Write([]byte(message))
		})

		if err := http.ListenAndServe(":"+strconv.Itoa(httpPort), nil); err != nil {
			panic(err)
		}
	}()

	var mutex = &sync.RWMutex{}
	storage := NewStorage(nil, 2000*time.Millisecond, mutex)
	id, err := storage.Register("webserver", "localhost", httpPort, []string{"http", "web"}, nil)

	if err != nil {
		t.Error(err)
	}

	err = storage.SetupHealthcheck(id, func() (bool, error) {
		service, err := storage.Service(&id, nil)
		conn, err := net.Dial("tcp", service.Address)
		if err != nil {
			return false, fmt.Errorf("connection error: %s", err.Error())
		}
		defer conn.Close()
		return true, nil
	})
	if err != nil {
		t.Error(err)
	}

	err = storage.Healthcheck(mutex)
	if err != nil {
		t.Error(err)
	}

	service, err := storage.Service(&id, nil)
	if err != nil {
		t.Error(err)
	}

	if service.Address != "localhost:8082" {
		t.Errorf("Address shoudl be localhost:8082, instead of %s", service.Address)
	}

	err = storage.Deregister(&id, nil)
	if err != nil {
		t.Error(err)
	}

}
