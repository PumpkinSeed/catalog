package consulmock

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/PumpkinSeed/catalog"
	"github.com/PumpkinSeed/consul/api"
)

const delimiter = "\n"
const delimiterByte = '\n'

type Catalog interface {
	Register(req *api.CatalogRegistration, q *api.WriteOptions) (*api.WriteMeta, error)
	Deregister(dereg *api.CatalogDeregistration, q *api.WriteOptions) (*api.WriteMeta, error)
	Datacenters() ([]string, error)
	Nodes(q *api.QueryOptions) ([]*api.Node, *api.QueryMeta, error)
	Node(node string, q *api.QueryOptions) (*api.CatalogNode, *api.QueryMeta, error)
	Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error)
	Service(service, tag string, q *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error)
}

type catalogmock struct {
	addr string
}

func NewCatalog(addr string) Catalog {
	return &catalogmock{addr: addr}
}

func (a *catalogmock) Register(req *api.CatalogRegistration, q *api.WriteOptions) (*api.WriteMeta, error) {
	var timeMeasure = time.Now()

	var rr = catalog.RegisterRequest{}
	a.translateRegisterRequest(req, &rr)
	rrJSON, err := json.Marshal(rr)
	if err != nil {
		return nil, err
	}

	var mainRequest = catalog.Request{
		Cmd: catalog.Register,
		Req: string(rrJSON),
	}

	resp, err := a.do(mainRequest)
	if err != nil {
		return nil, err
	}

	var respRegister catalog.RegisterResponse
	err = json.Unmarshal([]byte(resp.Resp), &respRegister)
	if err != nil {
		return nil, err
	}

	if respRegister.Success {
		return &api.WriteMeta{
			RequestTime: time.Since(timeMeasure),
		}, nil
	}
	return nil, errors.New(respRegister.Error)
}

func (a *catalogmock) Deregister(dereg *api.CatalogDeregistration, q *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, nil
}

func (a *catalogmock) Datacenters() ([]string, error) {
	return nil, nil
}

func (a *catalogmock) Nodes(q *api.QueryOptions) ([]*api.Node, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (a *catalogmock) Node(node string, q *api.QueryOptions) (*api.CatalogNode, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (a *catalogmock) Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (a *catalogmock) Service(service, tag string, q *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (a *catalogmock) translateRegisterRequest(req *api.CatalogRegistration, rr *catalog.RegisterRequest) error {
	rr.Name = req.Service.ID
	rr.Address = req.Service.Address
	rr.Port = req.Service.Port
	rr.Tags = req.Service.Tags

	return nil
}

func (a *catalogmock) translateDeregisterRequest(req *api.CatalogDeregistration, dr *catalog.DeregisterRequest) error {
	id := catalog.Identifier(12)
	dr.Name = &req.Address // @TODO ???
	dr.ID = &id            // @TODO ????

	return nil
}

func (a *catalogmock) do(req catalog.Request) (*catalog.Response, error) {
	rJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", a.addr)
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
