package output

import (
	"fmt"
	"os"
	"time"

	"github.com/nashabanov/urlcheck/internal/types"
)

type Config struct {
	ColorOutput bool
}

type Writer struct {
	config Config
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
)

func NewWriter(config Config) *Writer {
	return &Writer{config: config}
}

func (w *Writer) WriteProgress(current, total int, result types.Result) {
	var status string
	if result.Error != nil {
		status = w.colorize("✗", ColorRed)
	} else if result.StatusCode >= 200 && result.StatusCode < 300 {
		status = w.colorize("✓", ColorGreen)
	} else {
		status = w.colorize("!", ColorYellow)
	}

	var details string
	if result.Error != nil {
		details = fmt.Sprintf("(%s)", result.Error)
	} else {
		details = fmt.Sprintf("(%d, %v)", result.StatusCode, result.Duration)
	}

	totalDigits := len(fmt.Sprintf("%d", total))
	progress := fmt.Sprintf("[%*d/%d]", totalDigits, current, total)

	fmt.Printf("%s %s %s %s\n", progress, status, result.URL, details)
}

func (w *Writer) colorize(text, color string) string {
	if !w.config.ColorOutput || !isColorSupported() {
		return text
	}
	return color + text + ColorReset
}

func isColorSupported() bool {
	if os.Getenv("TERM") == "" {
		return false
	}
	fileinfo, _ := os.Stdout.Stat()
	return (fileinfo.Mode() & os.ModeCharDevice) != 0
}

type Summary struct {
	Total    int
	Success  int
	Failed   int
	Duration time.Duration
}

func (w *Writer) WriteSummary(summary Summary) {
	fmt.Println()

	successRate := float64(summary.Success) / float64(summary.Total) * 100

	successText := w.colorize(fmt.Sprintf("%d successful", summary.Success), ColorGreen)
	failedText := w.colorize(fmt.Sprintf("%d failed", summary.Failed), ColorRed)

	fmt.Printf("Summary: %s, %s, %.1f%% success rate\n", successText, failedText, successRate)
	fmt.Printf("Total: %d URLs checked in %v\n", summary.Total, summary.Duration.Round(time.Millisecond))
}
