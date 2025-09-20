package checker

import (
	"net"
	"net/http"
	"time"

	"github.com/nashabanov/urlcheck/internal/types"
)

type HTTPChecker struct {
	Timeout time.Duration
}

func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		Timeout: 5 * time.Second,
	}
}

func (hc *HTTPChecker) Check(url string) *types.Result {
	start := time.Now()
	client := &http.Client{Timeout: hc.Timeout}
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		var typedErr error

		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				typedErr = ErrTimeout{URL: url}
			} else if netErr.Temporary() {
				typedErr = ErrConnectionRefused{URL: url}
			} else {
				typedErr = ErrNetwork{URL: url}
			}
		} else if _, ok := err.(*net.DNSError); ok {
			typedErr = ErrDNSFailed{URL: url}
		} else {
			typedErr = ErrNetwork{URL: url}
		}

		return &types.Result{
			URL:        url,
			StatusCode: 0,
			Duration:   duration,
			Error:      typedErr,
		}
	}

	defer resp.Body.Close()
	return &types.Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}
