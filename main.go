package main

import (
	"context"
	"log"
	"time"
	"urlcheck/checker"
	"urlcheck/input"
	"urlcheck/output"
	"urlcheck/types"
	"urlcheck/worker"
)

func main() {
	// Input
	cfg := input.NewConfig("", []string{"https://google.com", "https://badurl123.com", "https://httpbin.org/status/404"}, false, 10)
	urls, err := cfg.GetURLs()
	if err != nil {
		log.Fatal(err)
	}

	// Components
	worker := &worker.Worker{MaxWorkers: 5}
	http_checker := checker.NewHTTPChecker()
	outputWriter := output.NewWriter(output.Config{ColorOutput: true})

	// Сбор результатов для summary
	var allResults []types.Result
	startTime := time.Now()

	// Run с накоплением
	ctx := context.Background()
	err = worker.Run(ctx, http_checker, urls, func(current, total int, result *types.Result) {
		outputWriter.WriteProgress(current, total, *result)
		allResults = append(allResults, *result) // ← Ключевая строка
	})

	if err != nil {
		log.Fatal(err)
	}

	// Summary
	duration := time.Since(startTime)
	summary := calculateSummary(allResults, duration)
	outputWriter.WriteSummmary(summary)
}

func calculateSummary(results []types.Result, duration time.Duration) output.Summary {
	total := len(results)
	success := 0

	for _, result := range results {
		if result.Error == nil && result.StatusCode >= 200 && result.StatusCode < 300 {
			success++
		}
	}

	return output.Summary{
		Total:    total,
		Success:  success,
		Failed:   total - success,
		Duration: duration,
	}
}
