package xhttp

import (
	"Predator/pkg/config"
	"Predator/pkg/xtls"
	"context"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/kataras/golog"
	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type (
	// RequestMiddleware type is for request middleware, called before a request is sent
	RequestMiddleware func(*Request, *Client) error
	// ResponseMiddleware type is for response middleware, called after a response has been received
	ResponseMiddleware func(*Response, *Client) error
)

type Client struct {
	HTTPClient    *http.Client
	ClientOptions *ClientOptions
	Error         interface{} // todo error handle exp

	defaultBeforeRequest []RequestMiddleware
	extraBeforeRequest   []RequestMiddleware
	afterResponse        []ResponseMiddleware
}

var defaultClient *Client
var defaultRedirectClient *Client
var requestCache gcache.Cache

func Init() error {
	// set no redirect client
	if DefaultClient() == nil {
		client, err := NewClient(false)
		if err != nil {
			return err
		}
		setDefaultClient(client)
	}
	// set redirect client
	if DefaultRedirectClient() == nil {
		client, err := NewClient(true)
		if err != nil {
			return err
		}
		setDefaultRedirectClient(client)
	}

	httpClientOptions := GetHTTPOptions()
	if requestCache == nil {
		gc := gcache.New(httpClientOptions.CacheMaxNumber).LRU().Build()
		requestCache = gc
	}

	return nil
}

func setDefaultClient(c *Client) {
	if c != nil {
		defaultClient = c
	}
}

func DefaultClient() *Client {
	return defaultClient
}

func setDefaultRedirectClient(c *Client) {
	if c != nil {
		defaultRedirectClient = c
	}
}

func DefaultRedirectClient() *Client {
	return defaultRedirectClient
}

func DoExt(ctx context.Context, redirect bool, req *Request) (resp *Response, err error) {
	if redirect {
		return defaultRedirectClient.Do(ctx, req)
	} else {
		return defaultClient.Do(ctx, req)
	}
}

func Do(ctx context.Context, req *Request) (resp *Response, err error) {
	return defaultClient.Do(ctx, req)
}

func DoWithRedirect(ctx context.Context, req *Request) (resp *Response, err error) {
	return defaultRedirectClient.Do(ctx, req)
}

func NewClient(followRedirects bool) (*Client, error) {
	httpClientOptions := GetHTTPOptions()
	if config.PocNum != 25 {
		httpClientOptions.MaxIdleConnsPerHost = config.PocNum * 2
	}
	if config.WebTimeout != 3 {
		httpClientOptions.ReadTimeout = int(time.Duration(config.WebTimeout) * time.Second)
	}
	hc, err := createHttpClient(followRedirects, httpClientOptions)
	if err != nil {
		return nil, err
	}

	c := &Client{
		HTTPClient:    hc,
		ClientOptions: httpClientOptions,
	}
	c.extraBeforeRequest = []RequestMiddleware{}
	c.defaultBeforeRequest = []RequestMiddleware{
		verifyRequestMethod,
		createHTTPRequest,
		readRequestBody,
	}
	c.afterResponse = []ResponseMiddleware{
		readResponseBody,
		responseLogger,
	}
	return c, nil
}

func createHttpClient(followRedirects bool, httpClientOptions *ClientOptions) (*http.Client, error) {
	tlsClientConfig, err := xtls.NewTLSConfig(xtls.DefaultClientOptions())
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(httpClientOptions.DialTimeout) * time.Second,
		}).DialContext,
		MaxConnsPerHost:       httpClientOptions.MaxConnsPerHost,
		MaxIdleConnsPerHost:   httpClientOptions.MaxIdleConnsPerHost,
		ResponseHeaderTimeout: time.Duration(httpClientOptions.ReadTimeout) * time.Second,
		IdleConnTimeout:       time.Duration(httpClientOptions.IdleConnTimeout) * time.Second,
		TLSHandshakeTimeout:   time.Duration(httpClientOptions.TLSHandshakeTimeout) * time.Second,
		MaxIdleConns:          httpClientOptions.MaxIdleConns,
		TLSClientConfig:       tlsClientConfig,
		DisableKeepAlives:     httpClientOptions.DisableKeepAlives,
	}
	if httpClientOptions.EnableHTTP2 {
		err := http2.ConfigureTransport(transport)
		if err != nil {
			return nil, err
		}
	}

	if httpClientOptions.Proxy != "" {
		proxy, err := url.Parse(httpClientOptions.Proxy)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Jar:           cookieJar,
		Transport:     transport,
		CheckRedirect: makeCheckRedirectFunc(followRedirects, httpClientOptions.MaxRedirect),
		Timeout:       time.Second * time.Duration(httpClientOptions.ReadTimeout),
	}, nil
}

type checkRedirectFunc func(req *http.Request, via []*http.Request) error

func makeCheckRedirectFunc(followRedirects bool, maxRedirects int) checkRedirectFunc {
	return func(req *http.Request, via []*http.Request) error {
		if !followRedirects {
			return http.ErrUseLastResponse
		}
		if len(via) >= maxRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

func (c *Client) Do(ctx context.Context, req *Request) (resp *Response, err error) {
	var (
		rawResp         *http.Response
		shouldRetry     bool
		doErr, retryErr error
	)

	if c == nil {
		golog.Debugf("xhttp client not instantiated. auto new no-follow-redirect client")
		tmpClient, err := NewClient(false)
		if err != nil {
			return nil, err
		}
		c = tmpClient
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(c.ClientOptions.MaxRequestTimeout)*time.Second)
	defer cancel()

	req.setContext(ctx)
	req.attempt = 0

	// user diy RequestMiddleware
	for _, f := range c.extraBeforeRequest {
		if err = f(req, c); err != nil {
			return nil, err
		}
	}

	// default diy RequestMiddleware
	for _, f := range c.defaultBeforeRequest {
		if err = f(req, c); err != nil {
			return nil, err
		}
	}

	var cacheErr error
	var reqMD5 string
	if req.GetHost() != "dnsbook.xyz" {
		reqMD5, cacheErr = req.MD5()
		if cacheErr == nil && c.ClientOptions.CacheMaxNumber > 0 && requestCache.Has(reqMD5) {
			meta, cacheErr := requestCache.GetIFPresent(reqMD5)
			if cacheErr == nil && meta != nil {
				cacheResponse, ok := meta.(**Response)
				if ok {
					golog.Debugf("%s md5 %s match the request cache, all match number %d", req.String(), reqMD5, requestCache.Len(false))
					return *cacheResponse, nil
				}
			}
		}
	}

	limiter := ExtractQPSLimiter(ctx)

	// do request with retry
	for i := 0; ; i++ {
		req.attempt++

		// qps limit
		err = limiter.Wait(req.GetContext(), req.RawRequest.Host)
		if err != nil {
			return nil, err
		}

		req.setSendAt()
		rawResp, doErr = c.HTTPClient.Do(req.RawRequest)
		// need retry
		shouldRetry, retryErr = defaultRetryPolicy(req.GetContext(), rawResp, doErr)
		if !shouldRetry {
			break
		}
		remain := c.ClientOptions.FailRetries - i
		if remain <= 0 {
			break
		}
		// waitTime
		waitTime := defaultBackoff(defaultRetryWaitMin, defaultRetryWaitMax, i, rawResp)
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("request over timeout %ds", c.ClientOptions.MaxRequestTimeout)
		case <-time.After(waitTime):
		case <-req.GetContext().Done():
			return nil, req.GetContext().Err()
		}
	}

	if doErr == nil && retryErr == nil && !shouldRetry {
		// request success
		response := &Response{
			Request:     req,
			RawResponse: rawResp,
		}
		response.setReceivedAt()

		// ResponseMiddleware
		for _, f := range c.afterResponse {
			if err = f(response, c); err != nil {
				return nil, err
			}
		}

		if req.GetHost() != "dnsbook.xyz" {
			if cacheErr == nil {
				err = requestCache.Set(reqMD5, &resp)
				if err != nil {
					return
				}
			}
		}

		return response, nil
	} else {
		finalErr := doErr
		if retryErr != nil {
			finalErr = retryErr
		}
		return nil, fmt.Errorf("giving up request to %s %s after %d attempt(s): %v",
			req.RawRequest.Method, req.RawRequest.URL, req.attempt, finalErr)
	}
}

func (c *Client) SetSkipVerifyTLS(skipVerify bool) {
	c.HTTPClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = skipVerify
}

func (c *Client) BeforeRequest(fn RequestMiddleware) {
	c.extraBeforeRequest = append(c.extraBeforeRequest, fn)
}

func (c *Client) AfterResponse(fn ResponseMiddleware) {
	c.afterResponse = append(c.afterResponse, fn)
}
