package catalog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

var binAddr = "127.0.0.1:7777"

func TestRegisterCommand(t *testing.T) {

	go func() {
		panic(NewServer(binAddr).Listen())
	}()

	var rr = RegisterRequest{
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

	tcpReq(t, req)
}

func tcpReq(t *testing.T, req Request) {
	rJSON, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}

	conn, err := net.Dial("tcp", binAddr)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Fprintf(conn, string(rJSON)+"\n")

	// listen for reply
	message, _ := bufio.NewReader(conn).ReadBytes(delimiterByte)
	var resp Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		t.Error(err)
		return
	}

	//var result RegisterResponse
	fmt.Printf("Message from server: %v", resp.Resp)
}
