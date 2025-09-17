package checker

import "time"

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      error
}

type Checker interface {
	Check(url string) *Result
}
