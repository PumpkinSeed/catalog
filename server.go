package catalog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type command uint8

const (
	Register command = iota
	Deregister
	Services
	Service
)

const delimiter = "\n"
const delimiterByte = '\n'

// Server represent the standalone service
type Server interface {
	Listen() error
	Close()
}

// server is the main handler struct
type server struct {
	bindAddr string
	storage  Storage
	//services map[identifier]ServiceSpec
	closeCh chan bool
}

func NewServer(bindAddr string) Server {
	closeCh := make(chan bool, 1)
	s := new(server)
	//s.services = make(map[identifier]ServiceSpec)
	s.storage = NewStorage()
	s.bindAddr = bindAddr
	s.closeCh = closeCh

	return s
}

func (s *server) Listen() error {
	ln, err := net.Listen("tcp", s.bindAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		select {
		case <-s.closeCh:
			fmt.Println("closed")
			return nil
		default:
		}

		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}

		fmt.Print("Message Received:", string(msg))

		resp, err := s.handleRequest(msg)
		if err != nil {
			fmt.Println(err)
			return err
		}

		n, err := conn.Write(resp)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Bytes written: %d\n", n)
	}

	return nil
}

func (s *server) Close() {
	s.closeCh <- true
	return
}

func (s *server) handleRequest(reqByte []byte) ([]byte, error) {
	var req Request
	var resp Response

	err := json.Unmarshal(reqByte, &req)
	if err != nil {
		return nil, err
	}

	switch req.Cmd {
	case Register:
		var registerReq RegisterRequest
		err = json.Unmarshal([]byte(req.Req), &registerReq)
		if err != nil {
			return nil, err
		}

		var registerResp RegisterResponse
		err := s.register(&registerReq, &registerResp)
		if err != nil {
			return nil, err
		}

		respJSON, err := json.Marshal(registerResp)
		if err != nil {
			return nil, err
		}
		resp = NewResponse(respJSON)
	case Deregister:
		var deregisterReq DeregisterRequest
		err = json.Unmarshal([]byte(req.Req), &deregisterReq)
		if err != nil {
			return nil, err
		}

		var deregisterResp DeregisterResponse
		err := s.deregister(&deregisterReq, &deregisterResp)
		if err != nil {
			return nil, err
		}

		respJSON, err := json.Marshal(deregisterResp)
		if err != nil {
			return nil, err
		}
		resp = NewResponse(respJSON)
	case Service:
		var serviceReq ServiceRequest
		err = json.Unmarshal([]byte(req.Req), &serviceReq)
		if err != nil {
			return nil, err
		}

		var serviceResp ServiceResponse
		err := s.service(&serviceReq, &serviceResp)
		if err != nil {
			return nil, err
		}

		respJSON, err := json.Marshal(serviceResp)
		if err != nil {
			return nil, err
		}
		resp = NewResponse(respJSON)
	case Services:
		var servicesReq ServicesRequest
		err = json.Unmarshal([]byte(req.Req), &servicesReq)
		if err != nil {
			return nil, err
		}

		var servicesResp ServicesResponse
		err := s.services(&servicesReq, &servicesResp)
		if err != nil {
			return nil, err
		}

		respJSON, err := json.Marshal(servicesResp)
		if err != nil {
			return nil, err
		}
		resp = NewResponse(respJSON)
	}

	return resp.prepare(), nil
}

func (s *server) register(req *RegisterRequest, resp *RegisterResponse) error {
	id, err := s.storage.Register(req.Name, req.Address, req.Port, req.Tags, req.Additional)
	resp.Meta = *req
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		return err
	}

	resp.ID = id
	resp.Success = true
	return err
}

func (s *server) deregister(req *DeregisterRequest, resp *DeregisterResponse) error {
	err := s.storage.Deregister(req.ID)
	resp.Meta = *req
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		return err
	}

	resp.Success = true
	return nil
}

func (s *server) service(req *ServiceRequest, resp *ServiceResponse) error {
	var ss ServiceSpec
	var err error

	// ID first manner
	if req.ID != nil {
		ss, err = s.storage.Service(req.ID, nil)
	} else if req.Name != nil {
		ss, err = s.storage.Service(nil, req.Name)
	} else {
		resp.Error = ErrServiceRequestInvalid.Error()
		resp.Success = false
		return ErrServiceRequestInvalid
	}

	resp.Meta = *req
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		return err
	}

	resp.Service = ss
	resp.Success = true
	return nil
}

func (s *server) services(req *ServicesRequest, resp *ServicesResponse) error {
	specs := s.storage.Services()
	var container []ServiceSpec
	for _, spec := range specs {
		container = append(container, spec)
	}
	resp.Meta = *req
	resp.Success = true
	resp.Services = container

	return nil
}
