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

// Server represent the standalone service
type Server interface {
	Listen() error
	Close()
}

// server is the main handler struct
type server struct {
	bindAddr string
	services map[identifier]ServiceSpec
	closeCh  chan bool
}

func NewServer(bindAddr string) Server {
	closeCh := make(chan bool, 1)
	s := new(server)
	s.services = make(map[identifier]ServiceSpec)
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

		registerResp, err := s.register(&registerReq)
		if err != nil {
			return nil, err
		}

		resp.Resp = registerResp
	case Deregister:
		var deregisterReq DeregisterRequest
		err = json.Unmarshal([]byte(req.Req), &deregisterReq)
		if err != nil {
			return nil, err
		}

		deregisterResp, err := s.deregister(&deregisterReq)
		if err != nil {
			return nil, err
		}

		resp.Resp = deregisterResp
	case Service:
	case Services:
	}

	return resp.prepare(), nil
}

func (s *server) register(req *RegisterRequest) (RegisterResponse, error) {
	return RegisterResponse{}, nil
}

func (s *server) deregister(req *DeregisterRequest) (DeregisterResponse, error) {
	return DeregisterResponse{}, nil
}
