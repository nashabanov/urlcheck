package checker

import (
	"time"
	"urlcheck/internal/types"
)

type MockChecker struct {
	Delay time.Duration
}

func (mc *MockChecker) Check(url string) *types.Result {
	if mc.Delay > 0 {
		time.Sleep(mc.Delay)
	}
	return &types.Result{
		URL:        url,
		StatusCode: 200,
		Duration:   100 * time.Millisecond,
		Error:      nil,
	}
}
