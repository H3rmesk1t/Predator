package lib

import (
	"Predator/module/v1/ymlpoc/structs"
	"Predator/pkg/xhttp"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Client           *xhttp.Client
	ClientNoRedirect *xhttp.Client

	urlTypePool = sync.Pool{
		New: func() interface{} {
			return new(structs.UrlType)
		},
	}
	connectInfoTypePool = sync.Pool{
		New: func() interface{} {
			return new(structs.ConnInfoType)
		},
	}
	addrTypePool = sync.Pool{
		New: func() interface{} {
			return new(structs.AddrType)
		},
	}
	tracePool = sync.Pool{
		New: func() interface{} {
			return new(httptrace.ClientTrace)
		},
	}

	requestPool = sync.Pool{
		New: func() interface{} {
			return new(structs.Request)
		},
	}
	responsePool = sync.Pool{
		New: func() interface{} {
			return new(structs.Response)
		},
	}
	httpBodyBufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	httpBodyPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}
)

func DoRequest(req *xhttp.Request, redirect bool) (*xhttp.Response, int64, error) {
	var (
		milliseconds int64
		oResp        *xhttp.Response
		err          error
	)

	if req.Body == nil || len(req.Body) == 0 {
	} else {
		req.SetHeader("Content-Length", strconv.Itoa(req.GetContentLength()))
		if req.GetHeaders().Get("Content-Type") == "" {
			req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	start := time.Now()
	trace := tracePool.Get().(*httptrace.ClientTrace)
	trace.GotFirstResponseByte = func() {
		milliseconds = time.Since(start).Nanoseconds() / 1e6
	}

	req.RawRequest.WithContext(httptrace.WithClientTrace(req.RawRequest.Context(), trace))

	ctx := context.Background()
	if redirect {
		oResp, err = Client.Do(ctx, req)
	} else {
		oResp, err = ClientNoRedirect.Do(ctx, req)
	}

	if err != nil {
		return nil, 0, err
	}

	return oResp, milliseconds, nil
}

func GetRespBody(oResp *xhttp.Response) ([]byte, error) {
	body := httpBodyPool.Get().([]byte)
	defer httpBodyPool.Put(body)

	if oResp.GetHeaders().Get("Content-Encoding") == "gzip" {
		gr, _ := gzip.NewReader(bytes.NewReader(oResp.Body))
		defer gr.Close()
		for {
			buf := httpBodyBufPool.Get().([]byte)
			n, err := gr.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
			body = append(body, buf...)
			httpBodyBufPool.Put(buf)
		}
	} else {
		raw, err := ioutil.ReadAll(bytes.NewReader(oResp.Body))
		if err != nil {
			return nil, err
		}
		body = raw
	}
	return body, nil
}

func ParseUrl(u *url.URL) *structs.UrlType {
	urlType := urlTypePool.Get().(*structs.UrlType)

	urlType.Scheme = u.Scheme
	urlType.Domain = u.Hostname()
	urlType.Host = u.Host
	urlType.Port = u.Port()
	urlType.Path = u.Path
	urlType.Query = u.RawQuery
	urlType.Fragment = u.Fragment

	return urlType
}

func ParseHttpRequest(oReq *xhttp.Request) (*structs.Request, error) {
	var (
		req = requestPool.Get().(*structs.Request)
	)

	req.Method = oReq.GetMethod()
	req.Url = ParseUrl(oReq.GetUrl())

	headers := make(map[string]string)
	for k := range oReq.GetHeaders() {
		headers[k] = oReq.GetHeaders().Get(k)
	}
	req.Headers = headers

	req.ContentType = oReq.GetHeaders().Get("Content-Type")
	if oReq.Body != nil && req.Body != nil {
		req.Body = make([]byte, len(oReq.Body))
		copy(req.Body, oReq.Body)
	}

	return req, nil
}

func ParseHttpResponse(oResp *xhttp.Response, milliseconds int64) (*structs.Response, error) {
	var (
		resp             = responsePool.Get().(*structs.Response)
		err              error
		header           string
		rawHeaderBuilder strings.Builder
	)

	headers := make(map[string]string)
	resp.Status = int32(oResp.GetStatus())
	resp.Url = ParseUrl(oResp.Request.GetUrl())

	for k := range oResp.GetHeaders() {
		header = oResp.GetHeaders().Get(k)
		headers[k] = header

		rawHeaderBuilder.WriteString(k)
		rawHeaderBuilder.WriteString(": ")
		rawHeaderBuilder.WriteString(header)
		rawHeaderBuilder.WriteString("\n")
	}
	resp.Headers = headers
	resp.ContentType = oResp.GetHeaders().Get("Content-Type")
	// 原始请求头
	resp.RawHeader = []byte(strings.Trim(rawHeaderBuilder.String(), "\n"))

	// 原始http响应
	resp.Raw, err = httputil.DumpResponse(oResp.RawResponse, true)
	body, err := GetRespBody(oResp)
	if err != nil {
		return nil, err
	}

	if err != nil {
		resp.Raw = body
	}
	// http响应体
	resp.Body = body

	// 响应时间
	resp.Latency = milliseconds

	return resp, nil
}

func ParseTCPUDPRequest(content []byte) (*structs.Request, error) {
	var (
		req = requestPool.Get().(*structs.Request)
	)

	req.Raw = content

	return req, nil
}

func ParseTCPUDPResponse(content []byte, socket *net.Conn, transport string) (*structs.Response, error) {
	var (
		resp       = responsePool.Get().(*structs.Response)
		conn       = connectInfoTypePool.Get().(*structs.ConnInfoType)
		connection = *socket

		addr     string
		addrType *structs.AddrType
		addrList []string
		port     string
	)

	resp.Raw = content

	// source
	addr = connection.LocalAddr().String()
	addrList = strings.SplitN(addr, ":", 2)
	if len(addrList) == 2 {
		port = addrList[1]
	} else {
		port = ""
	}

	addrType = addrTypePool.Get().(*structs.AddrType)
	addrType.Transport = transport
	addrType.Addr = addr
	addrType.Port = port
	conn.Source = addrType

	// destination
	addr = connection.RemoteAddr().String()
	addrList = strings.SplitN(addr, ":", 2)
	if len(addrList) == 2 {
		port = addrList[1]
	} else {
		port = ""
	}

	addrType = addrTypePool.Get().(*structs.AddrType)
	addrType.Transport = transport
	addrType.Addr = addr
	addrType.Port = port
	conn.Source = addrType
	conn.Destination = addrType

	resp.Conn = conn

	return resp, nil
}

func PutUrlType(urlType *structs.UrlType) {
	urlTypePool.Put(urlType)
}

func PutConnectInfo(connInfo *structs.ConnInfoType) {
	connectInfoTypePool.Put(connInfo)
}

func PutAddrType(addrType *structs.AddrType) {
	addrTypePool.Put(addrType)
}

func PutRequest(request *structs.Request) {
	requestPool.Put(request)
}
func PutResponse(response *structs.Response) {
	responsePool.Put(response)
}
