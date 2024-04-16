package xhttp

import (
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func verifyRequestMethod(req *Request, c *Client) error {
	if req.RawRequest == nil {
		return errors.New("req.RawRequest is nil")
	}
	currentMethod := req.RawRequest.Method
	if funk.Contains(c.ClientOptions.AllowMethods, currentMethod) == false {
		return fmt.Errorf(`http method %s not allowed`, currentMethod)
	}
	return nil
}

func readRequestBody(req *Request, c *Client) error {
	_, err := req.GetBody()
	if err != nil {
		return err
	}
	return nil
}

func setTrace(req *Request) {
	if req.clientTrace == nil && req.trace {
		req.clientTrace = &clientTrace{}
	}
	if req.clientTrace != nil {
		req.ctx = req.clientTrace.createContext(req.GetContext())
	}
}

func setRequestHeader(req *Request) {
	if req.RawRequest.Header.Get("Accept-Language") == "" {
		req.RawRequest.Header.Set("Accept-Language", "en")
	}
	if req.RawRequest.Header.Get("Accept") == "" {
		req.RawRequest.Header.Set("Accept", "*/*")
	}
}

func setContest(req *Request) {
	if req.GetContext() != nil {
		req.RawRequest = req.RawRequest.WithContext(req.GetContext())
	}
}

func createHTTPRequest(req *Request, c *Client) error {
	setTrace(req)
	setRequestHeader(req)
	setContest(req)

	req.RawRequest.Close = c.ClientOptions.DisableKeepAlives
	for key, value := range c.ClientOptions.Headers {
		if len(req.RawRequest.Header.Values(key)) > 0 {
			continue
		}
		req.RawRequest.Header.Set(key, value)
	}
	if c.ClientOptions.Cookies != nil {
		for k, v := range c.ClientOptions.Cookies {
			req.RawRequest.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}
	return nil
}

func readResponseBody(resp *Response, c *Client) error {
	lr := io.LimitReader(resp.RawResponse.Body, c.ClientOptions.MaxRespBodySize)
	bodyBytes, err := ioutil.ReadAll(lr)
	if err != nil {
		if bodyBytes != nil {
			return err
		} else {
			return err
		}
	}
	resp.Body = bodyBytes
	defer resp.RawResponse.Body.Close()
	return nil
}

func responseLogger(resp *Response, c *Client) error {
	if c.ClientOptions.Debug {
		req := resp.Request
		reqString, err := req.GetRaw()
		if err != nil {
			return err
		}

		respString, err := resp.GetRaw()
		if err != nil {
			return err
		}

		latency := resp.GetLatency()

		reqLog := "\n==============================================================================\n" +
			"--- REQUEST ---\n" +
			fmt.Sprintf("%s  %s  %s\n", req.GetMethod(), req.GetUrl().String(), req.RawRequest.Proto) +
			fmt.Sprintf("HOST   : %s\n", req.RawRequest.URL.Host) +
			fmt.Sprintf("RequestString:\n%s\n", reqString) +
			"------------------------------------------------------------------------------\n" +
			"--- RESPONSE ---\n" +
			fmt.Sprintf("STATUS       : %s\n", resp.RawResponse.Status) +
			fmt.Sprintf("PROTO        : %s\n", resp.RawResponse.Proto) +
			fmt.Sprintf("RECEIVED AT  : %v\n", resp.getReceivedAt().Format(time.RFC3339Nano)) +
			fmt.Sprintf("Attempt Num  : %d\n", req.attempt) +
			fmt.Sprintf("TIME DURATION: %v\n", latency) +
			fmt.Sprintf("HOST   : %s\n", req.RawRequest.URL.Host) +
			fmt.Sprintf("ResponseString:\n%s\n", respString) +
			"------------------------------------------------------------------------------\n"
		fmt.Println(reqLog)
	}
	return nil
}
