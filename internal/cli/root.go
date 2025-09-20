package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"urlcheck/internal/checker"
	"urlcheck/internal/input"
	"urlcheck/internal/output"
	"urlcheck/internal/types"
	"urlcheck/internal/worker"
)

func Execute() error {
	config, err := ParseFlags()
	if err != nil {
		return err
	}

	if err := config.Validate(); err != nil {
		if err.Error() == "version displayed" {
			os.Exit(0) // Нормальный выход для -version
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		return err
	}

	return runApplication(config)
}

func runApplication(config *Config) error {
	ctx, cancel := setupGracefulShutdown()
	defer cancel()

	urls, err := getURLs(config)
	if err != nil {
		return fmt.Errorf("failed to get URLs: %w", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs found to check")
	}

	return executeURLCheck(ctx, config, urls)
}

func getURLs(config *Config) ([]string, error) {
	var urls []string

	if config.URLs != "" {
		urls = parseURLString(config.URLs)
	}

	inputConfig := input.NewConfig(config.File, urls, config.Stdin, config.MaxUrls)
	return inputConfig.GetURLs()
}

func parseURLString(urlStr string) []string {
	urls := strings.Split(urlStr, ",")
	result := make([]string, 0, len(urls))

	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url != "" {
			result = append(result, url)
		}
	}

	return result
}

func setupGracefulShutdown() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Fprintf(os.Stderr, "\nReceived interrupt signal, shutting down gracefully...\n")
		cancel()
	}()

	return ctx, cancel
}

func executeURLCheck(ctx context.Context, config *Config, urls []string) error {
	// Создаем компоненты
	workerInstance := &worker.Worker{MaxWorkers: config.Workers}

	httpChecker := checker.NewHTTPChecker()
	httpChecker.Timeout = config.Timeout

	outputWriter := output.NewWriter(output.Config{
		ColorOutput: config.Color && !config.Quiet,
	})

	// Сбор результатов для итоговой статистики
	var allResults []types.Result
	startTime := time.Now()

	// Показываем начальное сообщение
	if !config.Quiet {
		fmt.Printf("Checking %d URLs with %d workers (timeout: %v)...\n",
			len(urls), config.Workers, config.Timeout)
	}

	// Выполняем проверку с callback'ом
	err := workerInstance.Run(ctx, httpChecker, urls, func(current, total int, result *types.Result) {
		if !config.Quiet {
			outputWriter.WriteProgress(current, total, *result)
		}
		allResults = append(allResults, *result)
	})

	// Обрабатываем ошибки выполнения
	if err != nil {
		if err == context.Canceled {
			fmt.Fprintf(os.Stderr, "Operation cancelled by user\n")
			os.Exit(130) // Стандартный exit code для SIGINT
		}
		return fmt.Errorf("execution failed: %w", err)
	}

	// Выводим итоговую статистику
	if !config.Quiet {
		duration := time.Since(startTime)
		summary := calculateSummary(allResults, duration)

		outputWriter.WriteSummary(summary)
	}

	// Определяем exit code
	exitCode := calculateExitCode(allResults)

	if exitCode != 0 {
		os.Exit(exitCode)
	}

	return nil
}

func calculateSummary(results []types.Result, duration time.Duration) output.Summary {
	total := len(results)
	success := 0

	for _, result := range results {
		// Считаем успешными только 2xx статус коды без ошибок
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

// calculateExitCode определяет код выхода программы
func calculateExitCode(results []types.Result) int {
	for _, result := range results {
		if result.Error != nil || result.StatusCode < 200 || result.StatusCode >= 300 {
			return 1 // Есть неудачные проверки
		}
	}
	return 0 // Все проверки успешны
}
