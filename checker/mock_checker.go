package checker

import "time"

type MockChecker struct {
	Delay time.Duration
}

func (mc *MockChecker) Check(url string) *Result {
	if mc.Delay > 0 {
		time.Sleep(mc.Delay)
	}
	return &Result{
		URL:        url,
		StatusCode: 200,
		Duration:   100 * time.Millisecond,
		Error:      nil,
	}
}
