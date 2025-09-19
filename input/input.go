package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	File    string
	URLs    []string
	Stdin   bool
	MaxUrls int
}

func NewConfig(file string, urls []string, stdin bool, maxUrls int) *Config {
	if maxUrls < 0 {
		maxUrls = 10000
	}

	if urls == nil {
		urls = []string{}
	}

	return &Config{
		File:    file,
		URLs:    urls,
		Stdin:   stdin,
		MaxUrls: maxUrls,
	}
}

func (cfg *Config) GetURLs() ([]string, error) {
	if len(cfg.URLs) > 0 {
		return cfg.URLs, nil
	}
	if cfg.File != "" {
		return readURLsFromFile(cfg.File, cfg.MaxUrls)
	}
	if cfg.Stdin {
		return readURLsFromStdin(cfg.MaxUrls)
	}
	return nil, fmt.Errorf("не указаны URL, файл или флаг --stdin")
}

func readURLsFromFile(filename string, limit int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла %q: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urls := readFilterLines(scanner, limit)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %q: %w", filename, err)
	}

	return urls, nil
}

func readURLsFromStdin(limit int) ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	urls := readFilterLines(scanner, limit)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения stdin: %w", err)
	}

	return urls, nil
}

func processLine(line string) (string, bool) {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") {
		return "", false
	}

	return line, true
}

func readFilterLines(s *bufio.Scanner, limit int) []string {
	var urls []string
	for s.Scan() {
		if line, valid := processLine(s.Text()); valid {
			if len(urls) == limit {
				break
			} else {
				urls = append(urls, line)
			}
		}
	}
	return urls
}
