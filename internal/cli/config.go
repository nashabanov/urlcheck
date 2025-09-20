package cli

import (
	"fmt"
	"time"
)

type Config struct {
	File  string
	URLs  string
	Stdin bool

	Workers int
	Timeout time.Duration
	MaxUrls int

	Color bool
	Quiet bool

	Version bool
}

const (
	AppVersion = "1.0.0"
	AppName    = "urlcheck"
)

func (c *Config) Validate() error {
	if c.Version {
		fmt.Printf("%s version %s", AppName, AppVersion)
		return fmt.Errorf("version displayed")
	}

	sources := 0
	if c.File != "" {
		sources++
	}
	if c.URLs != "" {
		sources++
	}
	if c.Stdin {
		sources++
	}

	if sources == 0 {
		return fmt.Errorf("no URL source specified. Use -file, -urls, or -stdin")
	}

	if sources > 1 {
		return fmt.Errorf("specify only one URL source (-file, -urls, or -stdin)")
	}

	if c.Workers <= 0 {
		return fmt.Errorf("")
	}

	if c.Workers > 1000 {
		return fmt.Errorf("")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("")
	}

	if c.MaxUrls <= 0 {
		return fmt.Errorf("")
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Workers: 5,
		Timeout: 5 * time.Second,
		MaxUrls: 10000,
		Color:   true,
		Quiet:   false,
	}
}
