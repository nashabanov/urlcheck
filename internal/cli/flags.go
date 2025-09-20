package cli

import (
	"flag"
	"fmt"
)

// ParseFlags парсит аргументы командной строки
func ParseFlags() (*Config, error) {
	config := DefaultConfig()

	// Определяем флаги
	flag.StringVar(&config.File, "file", config.File,
		"File containing URLs (one per line)")
	flag.StringVar(&config.URLs, "urls", config.URLs,
		"Comma-separated list of URLs")
	flag.BoolVar(&config.Stdin, "stdin", config.Stdin,
		"Read URLs from stdin")

	flag.IntVar(&config.Workers, "workers", config.Workers,
		"Number of concurrent workers")
	flag.DurationVar(&config.Timeout, "timeout", config.Timeout,
		"Request timeout (e.g., 5s, 1m)")
	flag.IntVar(&config.MaxUrls, "max-urls", config.MaxUrls,
		"Maximum number of URLs to process")

	flag.BoolVar(&config.Color, "color", config.Color,
		"Enable colored output")
	flag.BoolVar(&config.Quiet, "quiet", config.Quiet,
		"Quiet mode - show errors only")

	flag.BoolVar(&config.Version, "version", config.Version,
		"Show version information")

	// Кастомная справка
	flag.Usage = func() {
		printUsage()
	}

	// Парсим
	flag.Parse()

	return config, nil
}

// printUsage выводит справку по использованию
func printUsage() {
	fmt.Printf(`%s - Fast concurrent URL checker

Usage:
  %s [options]

Data Sources (choose exactly one):
  -file string       File containing URLs (one per line)
  -urls string       Comma-separated list of URLs
  -stdin             Read URLs from stdin

Options:
  -workers int       Number of concurrent workers (default: 10)
  -timeout duration  Request timeout, e.g. 5s, 1m (default: 5s)  
  -max-urls int      Maximum URLs to process (default: 10000)
  -color             Enable colored output (default: true)
  -quiet             Quiet mode - show errors only (default: false)
  -version           Show version information

Examples:
  # Check URLs from command line
  %s -urls "https://google.com,https://github.com"
  
  # Check URLs from file with 20 workers
  %s -file urls.txt -workers 20
  
  # Read from stdin with custom timeout
  cat urls.txt | %s -stdin -timeout 10s
  
  # Quiet mode without colors (good for scripts)
  %s -file urls.txt -color=false -quiet

Exit Codes:
  0  All URLs successful
  1  Some URLs failed or error occurred
  130 Interrupted by user (Ctrl+C)

`, AppName, AppName, AppName, AppName, AppName, AppName)
}
