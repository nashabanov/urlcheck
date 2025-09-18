package worker

import (
	"context"
	"fmt"
	"testing"
	"time"
	"urlcheck/checker"
)

func TestWorker_Run_Basic(t *testing.T) {
	w := Worker{
		MaxWorkers: 2,
	}
	urls := []string{"a", "b", "c"}

	mc := &checker.MockChecker{}

	results := w.Run(context.Background(), mc, urls)

	for _, res := range results {
		if res.StatusCode != 200 {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if res.Error != nil {
			t.Errorf("Expected no error, got: %v", res.Error)
		}
		if res.Duration <= 0 {
			t.Errorf("Expected duration > 0, got %v", res.Duration)
		}
	}
}

func TestWorker_Run_WithTimeout(t *testing.T) {
	w := Worker{
		MaxWorkers: 5,
	}

	urls := make([]string, 100)
	for i := range 100 {
		urls[i] = fmt.Sprintf("url-%d", i)
	}

	mc := &checker.MockChecker{Delay: 10 * time.Millisecond}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	results := w.Run(ctx, mc, urls)

	if len(results) == 100 {
		t.Errorf("Expected fewer than 100 results due to timeout, got all 100")
	}

	if len(results) == 0 {
		t.Errorf("Expected at least some results, got none")
	}

	for _, res := range results {
		if res.StatusCode != 200 {
			t.Errorf("Expected status code 200, got %d", res.StatusCode)
		}
		if res.Error != nil {
			t.Errorf("Expected no error, got: %v", res.Error)
		}
		if res.Duration <= 0 {
			t.Errorf("Expected duration > 0, got %v", res.Duration)
		}
	}

	t.Logf("Processed %d/%d URLs before timeout", len(results), len(urls))
}
