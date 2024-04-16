package xhttp

import (
	"context"
	"crypto/x509"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var (
	redirectsErrorRegex = regexp.MustCompile(`stopped after \d+ redirects\z`)
	schemeErrorRegex    = regexp.MustCompile(`unsupported protocol scheme`)
	defaultRetryWaitMin = 10 * time.Millisecond
	defaultRetryWaitMax = 50 * time.Millisecond
)

type checkRetry func(ctx context.Context, resp *http.Response, err error) (bool, error)

func defaultRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	return baseRetryPolicy(resp, err)
}

func baseRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			if redirectsErrorRegex.MatchString(v.Error()) {
				return false, v
			}

			if schemeErrorRegex.MatchString(v.Error()) {
				return false, v
			}

			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, v
			}
		}

		return true, nil
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	if resp.StatusCode == 0 {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	return false, nil
}

func defaultBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return time.Second * time.Duration(sleep)
				}
			}
		}
	}

	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}
