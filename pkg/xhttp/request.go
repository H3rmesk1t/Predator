package xhttp

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Request struct {
	RawRequest *http.Request
	Error      interface{}
	Body       []byte

	attempt     int
	ctx         context.Context
	raw         []byte
	trace       bool
	sendAt      time.Time
	clientTrace *clientTrace

	paramParse interface{}
}

func (r *Request) IsWithParam() bool {
	nQuery := r.GetClearQuery()
	nBody, err := r.GetClearBody()
	if err != nil {
		return false
	}
	if nQuery != "" || nBody != "" {
		return true
	}
	return false
}

func (r *Request) GetClearQuery() string {
	var queryKeys []string
	for k, v := range r.RawRequest.URL.Query() {
		for i := 0; i < len(v); i++ {
			queryKeys = append(queryKeys, k)
		}
	}
	sort.Strings(queryKeys)
	return strings.Join(queryKeys, "&")
}

func (r *Request) GetClearBody() (string, error) {
	var bodyKeys []string
	if r.GetContentType() == "application/json" {
		var body map[string]interface{}
		d := json.NewDecoder(bytes.NewReader(r.Body))
		err := d.Decode(&body)
		if err != nil {
			return "", err
		}
		for k := range body {
			bodyKeys = append(bodyKeys, k)
		}
	} else if r.GetContentType() == "application/x-www-form-urlencoded" {
		var query url.Values
		query, err := url.ParseQuery(string(r.Body))
		if err != nil {
			return "", err
		}
		for k, v := range query {
			for i := 0; i < len(v); i++ {
				bodyKeys = append(bodyKeys, k)
			}
		}
	}
	sort.Strings(bodyKeys)
	return strings.Join(bodyKeys, "&"), nil
}

func NewRequest(method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{
		RawRequest: req,
	}, nil
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s", r.GetMethod(), r.GetUrl().String())
}

func (r *Request) MD5() (string, error) {
	raw, err := r.GetRaw()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(raw)), nil
}

func (r *Request) Clone() *Request {
	return &Request{
		RawRequest: r.RawRequest.Clone(r.RawRequest.Context()),
		Error:      r.Error,
		Body:       r.Body,
	}
}

func (r *Request) GetContentLength() int {
	return r.GetContentLength()
}

func (r *Request) GetHost() string {
	return r.RawRequest.URL.Host
}

func (r *Request) GetContext() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

func (r *Request) GetAttempt() int {
	return r.attempt
}

func (r *Request) GetUrl() *url.URL {
	return r.RawRequest.URL
}

func (r *Request) GetMethod() string {
	return r.RawRequest.Method
}

func (r *Request) GetHeaders() http.Header {
	return r.RawRequest.Header
}

func (r *Request) GetContentType() string {
	return r.RawRequest.Header.Get("content-type")
}

func (r *Request) GetBody() ([]byte, error) {
	if r.Body != nil {
		return r.Body, nil
	}
	if r.RawRequest.Body == nil {
		return nil, nil
	}

	body, i, err := drainBody(r.RawRequest.Body)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	var dest io.Writer = &b
	_, err = io.Copy(dest, body)

	if b.Len() > 0 && r.RawRequest.Header.Get("Content-Length") == "" {
		r.RawRequest.Header.Set("Content-Length", fmt.Sprintf("%d", b.Len()))
		r.RawRequest.ContentLength = int64(b.Len())
	}

	r.Body = b.Bytes()
	r.RawRequest.Body = i

	return r.Body, nil
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func (r *Request) GetRaw() ([]byte, error) {
	reqHeaderRaw, err := httputil.DumpRequest(r.RawRequest, false)
	if err != nil {
		return nil, err
	}

	r.raw = append(reqHeaderRaw, r.Body...)
	return r.raw, nil
}

func (r *Request) getTraceInfo() TraceInfo {
	ct := r.clientTrace
	if ct == nil {
		return TraceInfo{}
	}
	ti := TraceInfo{
		DNSLookup:     ct.dnsDone.Sub(ct.dnsStart),
		TLSHandshake:  ct.tlsHandshakeDone.Sub(ct.tlsHandshakeStart),
		ServerTime:    ct.gotFirstResponseByte.Sub(ct.gotConn),
		IsConnReused:  ct.gotConnInfo.Reused,
		IsConnWasIdle: ct.gotConnInfo.WasIdle,
		ConnIdleTime:  ct.gotConnInfo.IdleTime,
	}

	if ct.gotConnInfo.Reused {
		ti.TotalTime = ct.endTime.Sub(ct.getConn)
	} else {
		ti.TotalTime = ct.endTime.Sub(ct.dnsStart)
	}

	if !ct.connectDone.IsZero() {
		ti.TCPConnTime = ct.connectDone.Sub(ct.dnsDone)
	}

	if !ct.gotConn.IsZero() {
		ti.ConnTime = ct.gotConn.Sub(ct.getConn)
	}

	if !ct.gotFirstResponseByte.IsZero() {
		ti.ResponseTime = ct.endTime.Sub(ct.gotFirstResponseByte)
	}

	if ct.gotConnInfo.Conn != nil {
		ti.RemoteAddr = ct.gotConnInfo.Conn.RemoteAddr()
	}

	return ti
}

func (r *Request) EnableTrace() *Request {
	r.trace = true
	return r
}

func (r *Request) setSendAt() *Request {
	r.sendAt = time.Now()
	return r
}

func (r *Request) GetSchema() string {
	return r.GetUrl().Scheme
}

func (r *Request) setContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.RawRequest.Header.Set(key, value)
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

func (r *Request) SetHeaderMultiValues(headers map[string][]string) *Request {
	for key, values := range headers {
		r.SetHeader(key, strings.Join(values, ", "))
	}
	return r
}

func (r *Request) SetHeaderMulti(headers map[string]string) *Request {
	for key, value := range headers {
		r.SetHeader(key, value)
	}
	return r
}

func (r *Request) AddCookie(hc *http.Cookie) *Request {
	r.RawRequest.AddCookie(hc)
	return r
}

func (r *Request) SetBody(body []byte) *Request {
	r.RawRequest.Body = io.NopCloser(bytes.NewReader(body))
	return r
}

func (r *Request) SetBasicAuth(username, password string) *Request {
	r.RawRequest.SetBasicAuth(username, password)
	return r
}

func (r *Request) SetPath(path string) *Request {
	r.RawRequest.URL.Path = path
	return r
}

func (r *Request) GetPath() string {
	return r.RawRequest.URL.Path
}

func (r *Request) SetQueryParam(key, value string) *Request {
	params := r.RawRequest.URL.Query()
	params.Set(key, value)
	r.RawRequest.URL.RawQuery = params.Encode()
	return r
}

func (r *Request) AddQueryParam(key, value string) *Request {
	params := r.RawRequest.URL.Query()
	params.Add(key, value)
	r.RawRequest.URL.RawQuery = params.Encode()
	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.RawRequest.Method = method
	return r
}
