package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"unsafe"

	"github.com/PumpkinSeed/catalog"
)

const delimiter = "\n"
const delimiterByte = '\n'

type Catalog interface {
	Register(name string, host string, port int, tags []string, additional interface{}) (string, error)
	Deregister(id *string, name *string) error
	Service(id *string, name *string) (*catalog.ServiceSpec, error)
	Services() ([]catalog.ServiceSpec, error)
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
		return strconv.FormatUint(uint64(respRegister.ID), 10), nil
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
	var idUint uint64
	var sr catalog.ServiceRequest
	if id != nil {
		idUint, _ = strconv.ParseUint(*id, 10, 64)
		sr = catalog.ServiceRequest{
			ID: (*catalog.Identifier)(unsafe.Pointer(&idUint)),
		}
	} else if name != nil {
		sr = catalog.ServiceRequest{
			Name: name,
		}
	}

	srJSON, err := json.Marshal(sr)
	if err != nil {
		return nil, err
	}

	var mainRequest = catalog.Request{
		Cmd: catalog.Service,
		Req: string(srJSON),
	}

	resp, err := c.do(mainRequest)
	if err != nil {
		return nil, err
	}

	var respService catalog.ServiceResponse
	err = json.Unmarshal([]byte(resp.Resp), &respService)
	if err != nil {
		return nil, err
	}

	if respService.Success {
		return &respService.Service, nil
	}

	return nil, errors.New(respService.Error)
}
func (c *catalogapi) Services() ([]catalog.ServiceSpec, error) {
	var sr = catalog.ServicesRequest{}

	srJSON, err := json.Marshal(sr)
	if err != nil {
		return nil, err
	}

	var mainRequest = catalog.Request{
		Cmd: catalog.Services,
		Req: string(srJSON),
	}

	resp, err := c.do(mainRequest)
	if err != nil {
		return nil, err
	}

	var respServices catalog.ServicesResponse
	err = json.Unmarshal([]byte(resp.Resp), &respServices)
	if err != nil {
		return nil, err
	}

	if respServices.Success {
		return respServices.Services, nil
	}

	return nil, errors.New(respServices.Error)
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
