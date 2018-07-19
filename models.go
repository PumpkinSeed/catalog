package catalog

import "encoding/json"

type Request struct {
	Cmd command `json:"cmd"`
	Req string  `json:"req"`
}

type Response struct {
	Resp interface{} `json:"resp"`
}

// prepare the response
func (r *Response) prepare() []byte {
	resp, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return append(resp, []byte(delimiter)...)
}

func NewResponse(resp interface{}) Response {
	return Response{Resp: resp}
}

// RegisterRequest represent the register request to the server
type RegisterRequest struct {
	Address    string      `json:"address"`
	Port       int         `json:"port"`
	Tags       []string    `json:"tags"`
	Additional interface{} `json:"additional"`
}

// RegisterResponse represent the register response to the server
type RegisterResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	ID      identifier      `json:"id"`
	Meta    RegisterRequest `json:"meta"`
}

// prepare the response
func (r *RegisterResponse) prepare() []byte {
	resp, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return append(resp, []byte(delimiter)...)
}

// DeregisterRequest represent the deregister request to the server
type DeregisterRequest struct {
	ID identifier `json:"id"`
}

// DeregisterResponse represent the deregister response to the server
type DeregisterResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Meta    DeregisterRequest `json:"meta"`
}

// prepare the response
func (r *DeregisterResponse) prepare() []byte {
	resp, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return append(resp, []byte(delimiter)...)
}

type ServiceRequest struct {
	Address string   `json:"address"`
	Port    int      `json:"port"`
	Tags    []string `json:"tags"`
}

type ServiceResponse struct {
	Success bool           `json:"success"`
	Error   string         `json:"error"`
	Meta    ServiceRequest `json:"meta"`
	Service ServiceSpec    `json:"service"`
}

// prepare the response
func (r *ServiceResponse) prepare() []byte {
	resp, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return append(resp, []byte(delimiter)...)
}

type ServicesRequest struct {
	Address string   `json:"address"`
	Port    int      `json:"port"`
	Tags    []string `json:"tags"`
}

type ServicesResponse struct {
	Success  bool           `json:"success"`
	Error    string         `json:"error"`
	Meta     ServiceRequest `json:"meta"`
	Services []ServiceSpec  `json:"services"`
}

// prepare the response
func (r *ServicesResponse) prepare() []byte {
	resp, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return append(resp, []byte(delimiter)...)
}
