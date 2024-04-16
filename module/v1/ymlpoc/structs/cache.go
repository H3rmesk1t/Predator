package structs

import (
	"Predator/pkg/xhttp"
)

type HttpRequestCache struct {
	Request       *xhttp.Request
	ProtoRequest  *Request
	ProtoResponse *Response
}

type TCPUDPRequestCache struct {
	Response      []byte
	ProtoResponse *Response
}
