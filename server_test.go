package catalog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

var binAddr = "127.0.0.1:7777"
var serv Server
var testCounter = 0
var idOfNewService identifier

func init() {
	go func() {
		serv = NewServer(binAddr)
		panic(serv.Listen())
	}()
	go func() {
		for {
			if testCounter > 0 {
				serv.Close()
			}
		}
	}()
}

func TestRegisterCommand(t *testing.T) {
	var rr = RegisterRequest{
		Name:    "webserver",
		Address: "localhost",
		Port:    8080,
		Tags:    []string{"web", "http"},
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

	idOfNewService = respRegister.ID

	if respRegister.Success != true {
		t.Errorf("Operation success should be true, instead of %v", respRegister.Success)
	}

	testCounter++
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
