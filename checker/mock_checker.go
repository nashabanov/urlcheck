package checker

import "time"

type MockChecker struct{}

func (mc *MockChecker) Check(url string) *Result {
	return &Result{
		URL:        url,
		StatusCode: 200,
		Duration:   100 * time.Millisecond,
		Error:      nil,
	}
}
