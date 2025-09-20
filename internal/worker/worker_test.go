package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nashabanov/urlcheck/internal/checker"
	"github.com/nashabanov/urlcheck/internal/types"
)

func TestValidateMaxWorkers(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{-1, 1},
		{-999, 1},
		{1, 1},
		{5, 5},
		{100, 100},
	}

	for _, tc := range testCases {
		worker := &Worker{MaxWorkers: tc.input}
		worker.validateMaxWorkers()

		if worker.MaxWorkers != tc.expected {
			t.Errorf("Input: %d, expected: %d, got: %d",
				tc.input, tc.expected, worker.MaxWorkers)
		}
	}
}

func makeCollectorCallback() (func(int, int, *types.Result), *[]types.Result) {
	var results []types.Result

	callback := func(current, total int, result *types.Result) {
		results = append(results, *result)
	}

	return callback, &results
}

func TestWorker_BasicCallback(t *testing.T) {
	worker := &Worker{MaxWorkers: 1}
	callback, results := makeCollectorCallback()

	urls := []string{"http://example.com"}
	ctx := context.Background()

	err := worker.Run(ctx, &checker.MockChecker{}, urls, callback)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(*results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(*results))
	}
}

func TestWorker_ContextCancellation(t *testing.T) {
	worker := &Worker{MaxWorkers: 1}
	callback, _ := makeCollectorCallback()

	urls := []string{"http://slow1.com", "http://slow2.com", "http://slow3.com"}
	ctx, cancel := context.WithCancel(context.Background())

	slowChecker := &checker.MockChecker{Delay: 5 * time.Second}

	errChan := make(chan error, 1)

	go func() {
		err := worker.Run(ctx, slowChecker, urls, callback)
		errChan <- err
	}()

	time.Sleep(time.Millisecond * 100)
	cancel()

	select {
	case err := <-errChan:
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second * 3):
		t.Fatal("Worker didn't stop within timeout")
	}
}

func TestWorker_MaxWorkers(t *testing.T) {
	maxWorkers := 2
	worker := &Worker{MaxWorkers: maxWorkers}

	// Контролируемый checker с синхронизацией
	controlledChecker := &controlledMockChecker{
		activeCount: make(chan int, 10), // Буферизованный канал для счетчика
	}

	callback, _ := makeCollectorCallback()

	urls := []string{
		"http://test1.com", "http://test2.com", "http://test3.com",
		"http://test4.com", "http://test5.com",
	}

	ctx := context.Background()

	// Запускаем в горутине с обработкой ошибки
	errChan := make(chan error, 1)
	go func() {
		err := worker.Run(ctx, controlledChecker, urls, callback)
		errChan <- err
	}()

	// Проверяем что количество одновременных операций не превышает maxWorkers
	maxConcurrent := 0
	timeout := time.After(time.Second * 5)

	for {
		select {
		case activeCount := <-controlledChecker.activeCount:
			if activeCount > maxConcurrent {
				maxConcurrent = activeCount
			}
			if activeCount == 0 && maxConcurrent > 0 {
				// Все завершились, проверяем результат
				if maxConcurrent > maxWorkers {
					t.Errorf("Expected max %d concurrent workers, got %d", maxWorkers, maxConcurrent)
				}

				// Проверяем что worker.Run завершился без ошибки
				select {
				case err := <-errChan:
					if err != nil {
						t.Errorf("Worker.Run returned unexpected error: %v", err)
					}
				case <-time.After(time.Millisecond * 100):
					t.Error("Worker.Run didn't complete in time")
				}

				return
			}
		case <-timeout:
			t.Fatal("Test timeout")
		}
	}
}

// Контролируемый checker для отслеживания конкуренции
type controlledMockChecker struct {
	activeCount chan int
	mu          sync.Mutex
	active      int
}

func (c *controlledMockChecker) Check(url string) *types.Result {
	// Увеличиваем счетчик активных
	c.mu.Lock()
	c.active++
	currentActive := c.active
	c.mu.Unlock()

	// Отправляем текущий счетчик
	c.activeCount <- currentActive

	// Имитируем работу
	time.Sleep(time.Millisecond * 200)

	// Уменьшаем счетчик
	c.mu.Lock()
	c.active--
	newActive := c.active
	c.mu.Unlock()

	// Отправляем обновленный счетчик
	c.activeCount <- newActive

	return &types.Result{
		URL:        url,
		StatusCode: 200,
		Duration:   time.Millisecond * 200,
		Error:      nil,
	}
}

func TestWorker_Run_WithTimeout(t *testing.T) {
	worker := &Worker{
		MaxWorkers: 1,
	}
	callback, results := makeCollectorCallback()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	slowChecker := &checker.MockChecker{Delay: 2 * time.Second}

	urls := []string{"http://very-slow.com"}

	err := worker.Run(ctx, slowChecker, urls, callback)

	// Проверяем что получили timeout
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}

	// Callback мог быть не вызван из-за таймаута
	if len(*results) > 0 {
		t.Logf("Some results received before timeout: %d", len(*results))
	}
}
