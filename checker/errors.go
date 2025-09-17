package checker

import "fmt"

type ErrTimeout struct {
	URL string
}

func (e ErrTimeout) Error() string {
	return fmt.Sprintf("timeout for %s", e.URL)
}

type ErrDNSFailed struct {
	URL string
}

func (e ErrDNSFailed) Error() string {
	return fmt.Sprintf("DNS lookup failed for %s", e.URL)
}

type ErrConnectionRefused struct {
	URL string
}

func (e ErrConnectionRefused) Error() string {
	return fmt.Sprintf("connection refused for %s", e.URL)
}

type ErrNetwork struct {
	URL string
}

func (e ErrNetwork) Error() string {
	return fmt.Sprintf("network error for %s", e.URL)
}
