package input

import (
	"bufio"
	"os"
	"testing"
)

func TestNewConfig_Basic(t *testing.T) {
	filename := ""
	urls := []string{"a", "b", "c"}
	limit := 10

	cfg := NewConfig(filename, urls, false, limit)

	if cfg.File != filename {
		t.Errorf("Expected %s, got %s", filename, cfg.File)
	}
	for i, url := range cfg.URLs {
		if urls[i] != url {
			t.Errorf("Expected %s, got %s", urls[i], url)
		}
	}
	if cfg.Stdin {
		t.Errorf("Expected false, got true")
	}
	if cfg.MaxUrls != limit {
		t.Errorf("Expected %d, got %d", limit, cfg.MaxUrls)
	}
}

func TestNewConfig_DefaultLimit(t *testing.T) {
	filename := ""
	urls := []string{"a", "b", "c"}
	limit := -1
	expectedLimit := 10000

	cfg := NewConfig(filename, urls, false, limit)

	if cfg.File != filename {
		t.Errorf("Expected %s, got %s", filename, cfg.File)
	}
	for i, url := range cfg.URLs {
		if urls[i] != url {
			t.Errorf("Expected %s, got %s", urls[i], url)
		}
	}
	if cfg.Stdin {
		t.Errorf("Expected false, got true")
	}
	if cfg.MaxUrls != expectedLimit {
		t.Errorf("Expected %d, got %d", expectedLimit, cfg.MaxUrls)
	}
}

func TestNewConfig_DefaultUrls(t *testing.T) {
	filename := ""
	limit := 10

	cfg := NewConfig(filename, nil, false, limit)

	if cfg.File != filename {
		t.Errorf("Expected %s, got %s", filename, cfg.File)
	}
	if cfg.URLs == nil {
		t.Errorf("Expected []string, got nil")
	}
	if cfg.Stdin {
		t.Errorf("Expected false, got true")
	}
	if cfg.MaxUrls != limit {
		t.Errorf("Expected %d, got %d", limit, cfg.MaxUrls)
	}
}

func TestGetURLs_URLsOnly(t *testing.T) {
	filename := ""
	exp_urls := []string{"a", "b", "c"}
	limit := 10

	cfg := NewConfig(filename, exp_urls, false, limit)
	urls, err := cfg.GetURLs()

	if err != nil {
		t.Errorf("Expected nil as error, got %s", err)
	}
	for i, url := range urls {
		if exp_urls[i] != url {
			t.Errorf("Expected %s, got %s", exp_urls[i], url)
		}
	}
}

func TestGetURLs_FileOnly(t *testing.T) {
	content := "http://example.com\n# комментарий\nhttp://google.com\n"
	exp_urls := []string{"http://example.com", "http://google.com"}

	tmpFile, err := os.CreateTemp("", "test-urls-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	urls, err := NewConfig(tmpFile.Name(), nil, false, 10).GetURLs()
	if err != nil {
		t.Errorf("Expected nil as error, got %s", err)
	}

	for i, url := range urls {
		if exp_urls[i] != url {
			t.Errorf("Expected %s, got %s", exp_urls[i], url)
		}
	}
}

func TestGetURLs_StdinOnly(t *testing.T) {
	content := "http://example.com\n# комментарий\nhttp://google.com\n"
	exp_urls := []string{"http://example.com", "http://google.com"}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	tmpFile, err := os.CreateTemp("", "test-urls-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(content)
	tmpFile.Seek(0, 0)

	os.Stdin = tmpFile

	urls, err := NewConfig("", nil, true, 10).GetURLs()
	if err != nil {
		t.Errorf("Expected nil as error, got %s", err)
	}

	for i, url := range urls {
		if exp_urls[i] != url {
			t.Errorf("Expected %s, got %s", exp_urls[i], url)
		}
	}
}

func TestGetUrls_NoSource(t *testing.T) {
	urls, err := NewConfig("", nil, false, 0).GetURLs()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if urls != nil {
		t.Errorf("Expected nil URLs, got %v", urls)
	}
}

func TestGetUrls_NoFile(t *testing.T) {
	urls, err := NewConfig("test.t", nil, false, 0).GetURLs()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if urls != nil {
		t.Errorf("Expected nil URLs, got %v", urls)
	}
}

func TestGetUrls_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-urls-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	urls, err := NewConfig(tmpFile.Name(), nil, true, 0).GetURLs()

	if err != nil {
		t.Errorf("Expected nil, got '%s'", err)
	}

	if len(urls) > 0 {
		t.Errorf("Expected emty URLs slice, got %v", urls)
	}
}

func TestGetUrls_EmtyStdin(t *testing.T) {
	urls, err := NewConfig("", nil, true, 0).GetURLs()

	if err != nil {
		t.Errorf("Expected nil, got '%s'", err)
	}

	if len(urls) > 0 {
		t.Errorf("Expected emty URLs slice, got %v", urls)
	}
}

func TestReadFiltredLines_Limit(t *testing.T) {
	content := "http://example.com\n# комментарий\nhttp://google.com\nhttp://google1.com\nhttp://google2.com\n"

	tmpFile, err := os.CreateTemp("", "test-urls-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	urls := readFilterLines(scanner, 1)

	if len(urls) > 1 {
		t.Errorf("Expected slice with len 1, got len %d", len(urls))
	}
}

func TestProcessLine_EmptyString(t *testing.T) {
	line, valid := processLine("")

	if line != "" {
		t.Errorf("Expected emty line, got %s", line)
	}

	if valid {
		t.Errorf("Expected valid false, got true")
	}
}

func TestProcessLine_Comment(t *testing.T) {
	line, valid := processLine("#comment")

	if line != "" {
		t.Errorf("Expected emty line, got %s", line)
	}

	if valid {
		t.Errorf("Expected valid false, got true")
	}
}

func TestProcessLine_WithSpaces(t *testing.T) {
	line, valid := processLine("  http://example.com  ")

	if line != "http://example.com" {
		t.Errorf("Expected 'http://example.com', got %s", line)
	}

	if !valid {
		t.Errorf("Expected valid, got false")
	}
}
