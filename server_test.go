package catalog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"testing"
)

var binAddr = "127.0.0.1:8878"
var serv Server
var testCounter = 2
var idOfServices []Identifier
var mutex = &sync.RWMutex{}
var testServices = []*struct {
	id         Identifier
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

func init() {

	go func() {
		serv = NewServer(binAddr, nil, mutex)
		err := serv.Listen()
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		for {
			if testCounter > 0 {
				//serv.Close()
			}
		}
	}()
}

func TestRegisterCommand(t *testing.T) {
	for _, srv := range testServices {
		var rr = RegisterRequest{
			Name:    srv.name,
			Address: srv.host,
			Port:    srv.port,
			Tags:    srv.tags,
		}
		rrJSON, err := json.Marshal(rr)
		if err != nil {
			t.Error(err)
			return
		}

		var req = Request{
			Cmd: Register,
			Req: string(rrJSON),
		}

		resp := tcpReq(t, req)
		var respRegister RegisterResponse
		err = json.Unmarshal([]byte(resp.Resp), &respRegister)
		if err != nil {
			t.Error(err)
			return
		}

		srv.id = respRegister.ID
		srv.registered = true

		if respRegister.Success != true {
			t.Errorf("Operation success should be true, instead of %v", respRegister.Success)
		}

	}
	testCounter++
}

func TestServiceCommand(t *testing.T) {
	for _, srv := range testServices {
		var sr = ServiceRequest{
			ID: &srv.id,
		}

		srJSON, err := json.Marshal(sr)
		if err != nil {
			t.Error(err)
			return
		}

		var req = Request{
			Cmd: Service,
			Req: string(srJSON),
		}

		resp := tcpReq(t, req)
		var respService ServiceResponse
		err = json.Unmarshal([]byte(resp.Resp), &respService)
		if err != nil {
			t.Error(err)
			return
		}

		if respService.Service.Name != srv.name {
			t.Errorf("Name should be %s, instead of %s", srv.name, respService.Service.Name)
		}
		if respService.Service.Host != srv.host {
			t.Errorf("Host should be %s, instead of %s", srv.host, respService.Service.Host)
		}
		if respService.Service.Port != srv.port {
			t.Errorf("Port should be %d, instead of %d", srv.port, respService.Service.Port)
		}
		if respService.Success != true {
			t.Errorf("Operation success should be true, instead of %v", respService.Success)
		}
	}

	testCounter++
}

func TestDeregisterCommand(t *testing.T) {
	srv := testServices[0]
	var sr = DeregisterRequest{
		ID: &srv.id,
	}

	srJSON, err := json.Marshal(sr)
	if err != nil {
		t.Error(err)
		return
	}

	var req = Request{
		Cmd: Deregister,
		Req: string(srJSON),
	}

	resp := tcpReq(t, req)
	var respService DeregisterResponse
	err = json.Unmarshal([]byte(resp.Resp), &respService)
	if err != nil {
		t.Error(err)
		return
	}

	if respService.Success == true {
		srv.registered = false
	}
}

func TestServicesCommand(t *testing.T) {
	var sr = ServicesRequest{}

	srJSON, err := json.Marshal(sr)
	if err != nil {
		t.Error(err)
		return
	}

	var req = Request{
		Cmd: Services,
		Req: string(srJSON),
	}

	resp := tcpReq(t, req)
	var respServices ServicesResponse
	err = json.Unmarshal([]byte(resp.Resp), &respServices)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(respServices)
}

func tcpReq(t *testing.T, req Request) *Response {
	rJSON, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
		return nil
	}

	conn, err := net.Dial("tcp", binAddr)
	if err != nil {
		t.Error(err)
		return nil
	}
	defer conn.Close()

	fmt.Fprintf(conn, string(rJSON)+"\n")

	message, _ := bufio.NewReader(conn).ReadBytes(delimiterByte)
	var resp Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		t.Error(err)
		return nil
	}

	return &resp
}
