package types

import "time"

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      error
}