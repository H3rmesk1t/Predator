package xhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Response struct {
	Request     *Request
	RawResponse *http.Response
	Body        []byte

	raw        []byte
	receivedAt time.Time
}

func (r *Response) GetHeaders() http.Header {
	return r.RawResponse.Header
}

func (r *Response) GetLocation() (*url.URL, error) {
	return r.RawResponse.Location()
}

func (r *Response) GetCookies() []*http.Cookie {
	return r.RawResponse.Cookies()
}

func (r *Response) GetHeadersExt() map[string]string {
	headers := make(map[string]string)
	for k, vv := range r.GetHeaders() {
		if len(vv) > 0 {
			headers[k] = vv[0]
		}
	}

	if len(r.RawResponse.Cookies()) > 1 {
		headers["Set-Cookie"] = ""
		for _, c := range r.RawResponse.Cookies() {
			s := fmt.Sprintf("%s=%s", c.Name, c.Value)
			if sc := headers["Set-Cookie"]; sc != "" {
				headers["Set-Cookie"] = sc + "; " + s
			} else {
				headers["Set-Cookie"] = s
			}
		}
	}
	return headers
}

func (r *Response) GetContentType() string {
	return r.RawResponse.Header.Get("content-type")
}

func (r *Response) GetUrl() *url.URL {
	return r.RawResponse.Request.URL
}

func (r *Response) GetLatency() time.Duration {
	if r.Request.clientTrace != nil {
		return r.Request.getTraceInfo().TotalTime
	}
	return r.receivedAt.Sub(r.Request.sendAt)
}

func (r *Response) GetStatus() int {
	return r.RawResponse.StatusCode
}

func (r *Response) GetBody() []byte {
	return r.Body
}

func (r *Response) GetRaw() ([]byte, error) {
	respHeaderRaw, err := httputil.DumpResponse(r.RawResponse, false)
	if err != nil {
		return nil, err
	}

	r.raw = append(respHeaderRaw, r.Body...)
	return r.raw, nil
}

func (r *Response) getReceivedAt() time.Time {
	return r.receivedAt
}

func (r *Response) setReceivedAt() {
	r.receivedAt = time.Now()
	if r.Request.clientTrace != nil {
		r.Request.clientTrace.endTime = r.receivedAt
	}
}
