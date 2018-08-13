package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"unsafe"

	"github.com/PumpkinSeed/catalog"
)

const delimiter = "\n"
const delimiterByte = '\n'

type Catalog interface {
	Register(name string, host string, port int, tags []string, additional interface{}) (string, error)
	Deregister(id *string, name *string) error
	Service(id *string, name *string) (*catalog.ServiceSpec, error)
	Services() map[catalog.Identifier]*catalog.ServiceSpec
}

type catalogapi struct {
	addr string
}

func NewCatalog(addr string) Catalog {
	return &catalogapi{addr: addr}
}

func (c *catalogapi) Register(name string, host string, port int, tags []string, additional interface{}) (string, error) {
	var rr = catalog.RegisterRequest{
		Name:    name,
		Address: host,
		Port:    port,
		Tags:    tags,
	}

	rrJSON, err := json.Marshal(rr)
	if err != nil {
		return "", err
	}

	var mainRequest = catalog.Request{
		Cmd: catalog.Register,
		Req: string(rrJSON),
	}

	resp, err := c.do(mainRequest)
	if err != nil {
		return "", err
	}

	var respRegister catalog.RegisterResponse
	err = json.Unmarshal([]byte(resp.Resp), &respRegister)
	if err != nil {
		return "", err
	}

	if respRegister.Success {
		return string(respRegister.ID), nil
	}
	return "", errors.New(respRegister.Error)

}
func (c *catalogapi) Deregister(id *string, name *string) error {
	var dr = catalog.DeregisterRequest{
		ID:   (*catalog.Identifier)(unsafe.Pointer(&id)),
		Name: name,
	}

	drJSON, err := json.Marshal(dr)
	if err != nil {
		return err
	}

	var mainRequest = catalog.Request{
		Cmd: catalog.Deregister,
		Req: string(drJSON),
	}

	resp, err := c.do(mainRequest)
	if err != nil {
		return err
	}

	var respDeregister catalog.DeregisterResponse
	err = json.Unmarshal([]byte(resp.Resp), &respDeregister)
	if err != nil {
		return err
	}

	if respDeregister.Success {
		return nil
	}

	return errors.New(respDeregister.Error)
}
func (c *catalogapi) Service(id *string, name *string) (*catalog.ServiceSpec, error) {

	return nil, nil
}
func (c *catalogapi) Services() map[catalog.Identifier]*catalog.ServiceSpec {
	return nil
}

func (c *catalogapi) do(req catalog.Request) (*catalog.Response, error) {
	rJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fmt.Fprintf(conn, string(rJSON)+"\n")

	message, _ := bufio.NewReader(conn).ReadBytes(delimiterByte)
	var resp catalog.Response
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
