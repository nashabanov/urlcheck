package checker

import (
	"testing"
	"time"
)

func checkResult(t *testing.T, result *Result, expectedURL string) {
	if result.URL != expectedURL {
		t.Errorf("Expected URL to be '%s', got '%s'", expectedURL, result.URL)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	minDuration := 99 * time.Millisecond
	maxDuration := 101 * time.Millisecond
	if result.Duration < minDuration || result.Duration > maxDuration {
		t.Errorf("Expected duration between 99ms and 101ms, got %v", result.Duration)
	}

	if result.Error != nil {
		t.Errorf("Expected no error, but got: %v", result.Error)
	}
}

func TestMockChecker(t *testing.T) {
	url := "https://example.com"
	mc := MockChecker{}
	result := mc.Check(url)
	checkResult(t, result, url)
}

func TestMockChecker_EmptyUrl(t *testing.T) {
	url := ""
	mc := MockChecker{}
	result := mc.Check(url)
	checkResult(t, result, url)
}
